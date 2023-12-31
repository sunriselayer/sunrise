package keeper

import (
	"context"

	"cosmossdk.io/errors"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	gethcommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"sunrise/x/blobstream/types"
)

func (k msgServer) RegisterEvmAddress(goCtx context.Context, msg *types.MsgRegisterEvmAddress) (*types.MsgRegisterEvmAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	valAddr := sdk.ValAddress(addr)

	evmAddr := gethcommon.HexToAddress(msg.EvmAddress)

	if _, err := k.StakingKeeper.GetValidator(ctx, valAddr); err != nil {
		return nil, staking.ErrNoValidatorFound
	}

	if !k.IsEVMAddressUnique(ctx, evmAddr) {
		return nil, errors.Wrapf(types.ErrEVMAddressAlreadyExists, "address %s", msg.EvmAddress)
	}

	k.SetEVMAddress(ctx, valAddr, evmAddr)

	return &types.MsgRegisterEvmAddressResponse{}, nil
}
