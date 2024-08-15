package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetFaultCounter(ctx context.Context, operator sdk.ValAddress) uint64 {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.GetFaultCounterKey(operator))
	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetFaultCounter(ctx context.Context, operator sdk.ValAddress, faultCounter uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetFaultCounterKey(operator), sdk.Uint64ToBigEndian(faultCounter))
}

func (k Keeper) DeleteFaultCounter(ctx context.Context, operator sdk.ValAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.GetFaultCounterKey(operator))
}

func (k Keeper) IterateFaultCounters(ctx context.Context,
	handler func(operator sdk.ValAddress, faultCount uint64) (stop bool),
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(storeAdapter, types.FaultCounterKeyPrefix)
	iter := storetypes.KVStorePrefixIterator(prefixStore, []byte{})
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		operator := sdk.ValAddress(iter.Key())

		if handler(operator, sdk.BigEndianToUint64(iter.Value())) {
			break
		}
	}
}

func (k Keeper) HandleSlashEpoch(ctx sdk.Context) {
	params := k.GetParams(ctx)
	powerReduction := k.StakingKeeper.PowerReduction(ctx)
	k.IterateFaultCounters(ctx, func(operator sdk.ValAddress, faultCount uint64) bool {
		validator, err := k.StakingKeeper.Validator(ctx, operator)
		if err != nil {
			panic(err)
		}

		defer k.DeleteFaultCounter(ctx, operator)
		if validator.IsJailed() || !validator.IsBonded() {
			return false
		}

		if faultCount <= params.EpochMaxFault {
			return false
		}

		consAddr, err := validator.GetConsAddr()
		if err != nil {
			panic(err)
		}

		err = k.SlashingKeeper.Slash(
			ctx, consAddr, params.SlashFraction,
			validator.GetConsensusPower(powerReduction),
			ctx.BlockHeight()-sdk.ValidatorUpdateDelay-1,
		)
		if err != nil {
			panic(err)
		}
		err = k.SlashingKeeper.Jail(ctx, consAddr)
		if err != nil {
			panic(err)
		}
		return false
	})
}
