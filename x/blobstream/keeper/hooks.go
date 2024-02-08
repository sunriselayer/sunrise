package keeper

import (
	"context"

	"github.com/sunrise-zone/sunrise-app/x/blobstream/types"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Hooks is a wrapper struct around Keeper.
type Hooks struct {
	k Keeper
}

// Hooks Create new Blobstream hooks
func (k Keeper) Hooks() Hooks {
	// if startup is mis-ordered in app.go this hook will halt the chain when
	// called. Keep this check to make such a mistake obvious
	if k.storeService == nil {
		panic("hooks initialized before BlobstreamKeeper")
	}
	return Hooks{k}
}

func (h Hooks) AfterValidatorBeginUnbonding(ctx context.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// When Validator starts Unbonding, Persist the block height in the store
	// Later in endblocker, check if there is at least one validator who started
	// unbonding and create a valset request. The reason for creating valset
	// requests in endblock is to create only one valset request per block, if
	// multiple validators starts unbonding at same block.

	// this hook IS called for jailing or unbonding triggered by users but it IS
	// NOT called for jailing triggered in the endblocker therefore we call the
	// keeper function ourselves there.

	h.k.SetLatestUnBondingBlockHeight(sdkCtx, uint64(sdkCtx.BlockHeight()))
	return nil
}

func (h Hooks) BeforeDelegationCreated(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorCreated(ctx context.Context, addr sdk.ValAddress) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	defaultEvmAddr := types.DefaultEvmAddress(addr)
	// This should practically never happen that we have a collision. It may be
	// bad UX to reject the attempt to create a validator and require the user to
	// generate a new set of keys but this ensures EVM address uniqueness
	if !h.k.IsEVMAddressUnique(sdkCtx, defaultEvmAddr) {
		return errors.Wrapf(types.ErrEVMAddressAlreadyExists, "create a validator with a different operator address to %s (pubkey collision)", addr.String())
	}
	h.k.SetEVMAddress(sdkCtx, addr, defaultEvmAddr)
	return nil
}

func (h Hooks) BeforeValidatorModified(_ context.Context, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorBonded(_ context.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationRemoved(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorRemoved(_ context.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeValidatorSlashed(_ context.Context, _ sdk.ValAddress, _ math.LegacyDec) error {
	return nil
}

func (h Hooks) BeforeDelegationSharesModified(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterDelegationModified(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterUnbondingInitiated(_ context.Context, _ uint64) error {
	return nil
}
