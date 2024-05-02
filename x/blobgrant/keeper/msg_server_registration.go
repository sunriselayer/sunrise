package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise-app/x/blobgrant/types"
)

func (k msgServer) CreateRegistration(goCtx context.Context, msg *types.MsgCreateRegistration) (*types.MsgCreateRegistrationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetRegistration(
		ctx,
		msg.Address,
	)
	if isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var registration = types.Registration{
		Address:      msg.Address,
		ProxyAddress: msg.ProxyAddress,
	}

	k.SetRegistration(
		ctx,
		registration,
	)
	return &types.MsgCreateRegistrationResponse{}, nil
}

func (k msgServer) UpdateRegistration(goCtx context.Context, msg *types.MsgUpdateRegistration) (*types.MsgUpdateRegistrationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetRegistration(
		ctx,
		msg.Address,
	)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Address != valFound.Address {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var registration = types.Registration{
		Address:      msg.Address,
		ProxyAddress: msg.ProxyAddress,
	}

	k.SetRegistration(ctx, registration)

	return &types.MsgUpdateRegistrationResponse{}, nil
}

func (k msgServer) DeleteRegistration(goCtx context.Context, msg *types.MsgDeleteRegistration) (*types.MsgDeleteRegistrationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetRegistration(
		ctx,
		msg.Address,
	)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Address != valFound.Address {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveRegistration(
		ctx,
		msg.Address,
	)

	return &types.MsgDeleteRegistrationResponse{}, nil
}
