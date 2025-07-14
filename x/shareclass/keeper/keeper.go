package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
	logger       log.Logger

	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority string

	addressCodec address.Codec

	Schema                    collections.Schema
	Params                    collections.Item[types.Params]
	Unbondings                *collections.IndexedMap[uint64, types.Unbonding, types.UnbondingsIndexes]
	UnbondingId               collections.Sequence
	RewardMultiplier          collections.Map[collections.Pair[[]byte, string], string]                   // math.Dec
	UsersLastRewardMultiplier collections.Map[collections.Triple[sdk.AccAddress, []byte, string], string] // math.Dec
	LastRewardHandlingTime    collections.Map[[]byte, int64]

	accountKeeper        types.AccountKeeper
	bankKeeper           types.BankKeeper
	stakingKeeper        types.StakingKeeper
	tokenConverterKeeper types.TokenConverterKeeper
	distributionKeeper   types.DistributionKeeper

	StakingMsgServer   stakingtypes.MsgServer
	StakingQueryServer stakingtypes.QueryServer
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	addressCodec address.Codec,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distributionKeeper types.DistributionKeeper,
	tokenConverterKeeper types.TokenConverterKeeper,
	stakingMsgServer stakingtypes.MsgServer,
	stakingQueryServer stakingtypes.QueryServer,
) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		logger:       logger,
		authority:    authority,
		addressCodec: addressCodec,

		Params:                    collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Unbondings:                collections.NewIndexedMap(sb, types.UnbondingsKeyPrefix, "unbondings", types.UnbondingsKeyCodec, codec.CollValue[types.Unbonding](cdc), types.NewUnbondingsIndexes(sb, addressCodec)),
		UnbondingId:               collections.NewSequence(sb, types.UnbondingIdKey, "unbonding_id"),
		RewardMultiplier:          collections.NewMap(sb, types.RewardMultiplierKeyPrefix, "reward_multiplier", types.RewardMultiplierKeyCodec, collections.StringValue),
		UsersLastRewardMultiplier: collections.NewMap(sb, types.UsersLastRewardMultiplierKeyPrefix, "users_last_reward_multiplier", types.UsersLastRewardMultiplierKeyCodec, collections.StringValue),
		LastRewardHandlingTime:    collections.NewMap(sb, types.LastRewardHandlingTimeKeyPrefix, "last_reward_handling_time", types.LastRewardHandlingTimeKeyCodec, collections.Int64Value),

		accountKeeper:        accountKeeper,
		bankKeeper:           bankKeeper,
		stakingKeeper:        stakingKeeper,
		distributionKeeper:   distributionKeeper,
		tokenConverterKeeper: tokenConverterKeeper,

		StakingMsgServer:   stakingMsgServer,
		StakingQueryServer: stakingQueryServer,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
