package keeper

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashicorp/go-metrics"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/sunriselayer/sunrise/x/cron/types"
)

var (
	LabelExecuteReadySchedules   = "execute_ready_schedules"
	LabelExecuteCronSchedule     = "execute_cron_schedule"
	LabelExecuteCronContract     = "execute_cron_contract"
	LabelScheduleCount           = "schedule_count"
	LabelScheduleExecutionsCount = "schedule_executions_count"

	MetricLabelSuccess      = "success"
	MetricLabelScheduleName = "schedule_name"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte
	logger    log.Logger

	Schema        collections.Schema
	Params        collections.Item[types.Params]
	Schedules     collections.Map[string, types.Schedule]
	ScheduleCount collections.Item[types.ScheduleCount]

	accountKeeper types.AccountKeeper
	WasmMsgServer types.WasmMsgServer
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	logger log.Logger,
	accountKeeper types.AccountKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		logger:       logger,

		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Schedules:     collections.NewMap(sb, types.SchedulesKey, "schedules", collections.StringKey, codec.CollValue[types.Schedule](cdc)),
		ScheduleCount: collections.NewItem(sb, types.ScheduleCountKey, "schedule_count", codec.CollValue[types.ScheduleCount](cdc)),

		accountKeeper: accountKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) SetWasmMsgServer(server types.WasmMsgServer) {
	k.WasmMsgServer = server
}

// ExecuteReadySchedules gets all schedules that are due for execution (with limit that is equal to Params.Limit)
// and executes messages in each one
func (k Keeper) ExecuteReadySchedules(ctx sdk.Context, executionStage types.ExecutionStage) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), LabelExecuteReadySchedules)
	schedules, err := k.getSchedulesReadyForExecution(ctx, executionStage)
	if err != nil {
		k.Logger().Error("failed to get schedules ready for execution", "error", err)
		return
	}

	for _, schedule := range schedules {
		err := k.executeSchedule(ctx, schedule)
		recordExecutedSchedule(err, schedule)
	}
}

// AddSchedule adds a new schedule to be executed every certain number of blocks, specified in the `period`.
// First schedule execution is supposed to be on `now + period` block.
func (k Keeper) AddSchedule(
	ctx sdk.Context,
	name string,
	period uint64,
	msgs []types.MsgExecuteContract,
	executionStage types.ExecutionStage,
) error {
	exists, err := k.Schedules.Has(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("schedule already exists with name=%v", name)
	}

	schedule := types.Schedule{
		Name:              name,
		Period:            period,
		Msgs:              msgs,
		LastExecuteHeight: uint64(ctx.BlockHeight()),
		ExecutionStage:    executionStage,
	}

	if err := k.Schedules.Set(ctx, name, schedule); err != nil {
		return err
	}
	return k.changeTotalCount(ctx, 1)
}

// RemoveSchedule removes schedule with a given `name`
func (k Keeper) RemoveSchedule(ctx sdk.Context, name string) error {
	exists, err := k.Schedules.Has(ctx, name)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	if err := k.changeTotalCount(ctx, -1); err != nil {
		return err
	}
	return k.Schedules.Remove(ctx, name)
}

// GetSchedule returns schedule with a given `name`
func (k Keeper) GetSchedule(ctx sdk.Context, name string) (*types.Schedule, bool) {
	schedule, err := k.Schedules.Get(ctx, name)
	if err != nil {
		return nil, false
	}
	return &schedule, true
}

// GetAllSchedules returns all schedules
func (k Keeper) GetAllSchedules(ctx sdk.Context) ([]types.Schedule, error) {
	res := make([]types.Schedule, 0)
	err := k.Schedules.Walk(ctx, nil, func(_ string, schedule types.Schedule) (stop bool, err error) {
		res = append(res, schedule)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (k Keeper) GetScheduleCount(ctx sdk.Context) int32 {
	count, err := k.ScheduleCount.Get(ctx)
	if err != nil {
		return 0
	}
	return count.Count
}

func (k Keeper) getSchedulesReadyForExecution(ctx sdk.Context, executionStage types.ExecutionStage) ([]types.Schedule, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	count := uint64(0)

	res := make([]types.Schedule, 0)

	err = k.Schedules.Walk(ctx, nil, func(key string, schedule types.Schedule) (stop bool, err error) {
		if k.intervalPassed(ctx, schedule) && schedule.ExecutionStage == executionStage {
			res = append(res, schedule)
			count++

			if count >= params.Limit {
				k.Logger().Info("limit of schedule executions per block reached")
				return true, nil // stop iteration
			}
		}
		return false, nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

// executeSchedule executes all msgs in a given schedule and changes LastExecuteHeight
// if at least one msg execution fails, rollback all messages
func (k Keeper) executeSchedule(ctx sdk.Context, schedule types.Schedule) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), LabelExecuteCronSchedule, schedule.Name)
	schedule.LastExecuteHeight = uint64(ctx.BlockHeight())
	if err := k.Schedules.Set(ctx, schedule.Name, schedule); err != nil {
		return err
	}

	cacheCtx, writeFn := ctx.CacheContext()

	for idx, msg := range schedule.Msgs {
		startTimeContract := time.Now()
		executeMsg := wasmtypes.MsgExecuteContract{
			Sender:   k.accountKeeper.GetModuleAddress(types.ModuleName).String(),
			Contract: msg.Contract,
			Msg:      []byte(msg.Msg),
			Funds:    sdk.NewCoins(),
		}
		_, err := k.WasmMsgServer.ExecuteContract(sdk.WrapSDKContext(cacheCtx), &executeMsg)
		telemetry.ModuleMeasureSince(types.ModuleName, startTimeContract, LabelExecuteCronContract, schedule.Name, msg.Contract)
		if err != nil {
			k.Logger().Info("executeSchedule: failed to execute contract msg",
				"schedule_name", schedule.Name,
				"msg_idx", idx,
				"msg_contract", msg.Contract,
				"msg", msg.Msg,
				"error", err,
			)
			return err
		}
	}

	writeFn()
	return nil
}

func (k Keeper) intervalPassed(ctx sdk.Context, schedule types.Schedule) bool {
	return uint64(ctx.BlockHeight()) >= (schedule.LastExecuteHeight + schedule.Period)
}

func (k Keeper) changeTotalCount(ctx sdk.Context, incrementAmount int32) error {
	count, err := k.ScheduleCount.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			count = types.ScheduleCount{Count: 0}
		} else {
			return err
		}
	}

	newCount := types.ScheduleCount{Count: count.Count + incrementAmount}
	err = k.ScheduleCount.Set(ctx, newCount)
	if err != nil {
		return err
	}

	telemetry.ModuleSetGauge(types.ModuleName, float32(newCount.Count), LabelScheduleCount)
	return nil
}

func recordExecutedSchedule(err error, schedule types.Schedule) {
	telemetry.IncrCounterWithLabels([]string{LabelScheduleExecutionsCount}, 1, []metrics.Label{
		telemetry.NewLabel(telemetry.MetricLabelNameModule, types.ModuleName),
		telemetry.NewLabel(MetricLabelSuccess, strconv.FormatBool(err == nil)),
		telemetry.NewLabel(MetricLabelScheduleName, schedule.Name),
	})
}
