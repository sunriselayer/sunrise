package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/stable/types"
)

func (k msgServer) Burn(ctx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	authorityContract, err := k.addressCodec.StringToBytes(msg.AuthorityContract)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	amount := sdk.NewCoins(sdk.NewCoin(params.StableDenom, msg.Amount))

	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		authorityContract,
		types.ModuleName,
		amount,
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.BurnCoins(
		ctx,
		types.ModuleName,
		amount,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgBurnResponse{}, nil
}
