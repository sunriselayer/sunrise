// Package keeper provides the core logic for the stable module.
// This file implements the MsgServer interface for minting.
package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/stable/types"
)

// Mint handles the MsgMint request. It is a wrapper around the keeper's Mint function.
func (k msgServer) Mint(ctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "invalid sender address: %s", msg.Sender)
	}

	if err := msg.Amount.Validate(); err != nil {
		return nil, errorsmod.Wrap(err, "invalid collateral amount")
	}

	mintedAmount, err := k.Keeper.Mint(ctx, sender, msg.Amount)
	if err != nil {
		return nil, err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.MsgMintResponse{
		Amount: sdk.NewCoins(sdk.NewCoin(params.StableDenom, mintedAmount)),
	}, nil
}
