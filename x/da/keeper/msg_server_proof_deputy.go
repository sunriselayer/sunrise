package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RegisterProofDeputy(ctx context.Context, msg *types.MsgRegisterProofDeputy) (*types.MsgRegisterProofDeputyResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}
	deputy, err := k.addressCodec.StringToBytes(msg.DeputyAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid deputy address")
	}
	valAddr := sdk.ValAddress(sender)
	_, err = k.StakingKeeper.Validator(ctx, valAddr)
	if err != nil {
		return nil, err
	}
	err = k.SetProofDeputy(ctx, sender, deputy)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = sdkCtx.EventManager().EmitTypedEvent(msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgRegisterProofDeputyResponse{}, nil
}

func (k msgServer) UnregisterProofDeputy(ctx context.Context, msg *types.MsgUnregisterProofDeputy) (*types.MsgUnregisterProofDeputyResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}
	valAddr := sdk.ValAddress(sender)
	_, err = k.StakingKeeper.Validator(ctx, valAddr)
	if err != nil {
		return nil, err
	}
	_, found := k.GetProofDeputy(ctx, sender)
	if !found {
		return nil, types.ErrDeputyNotFound
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.DeleteProofDeputy(sdkCtx, sender)

	err = sdkCtx.EventManager().EmitTypedEvent(msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgUnregisterProofDeputyResponse{}, nil
}
