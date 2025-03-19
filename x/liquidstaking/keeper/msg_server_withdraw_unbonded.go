package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k msgServer) WithdrawUnbonded(ctx context.Context, msg *types.MsgWithdrawUnbonded) (*types.MsgWithdrawUnbondedResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	params, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: iterate Unstakings to calculate unbonded amount, and delete used Unstaking state

	// Convert bond denom to fee denom
	err = k.tokenConverterKeeper.Convert(ctx, amount, sender)
	if err != nil {
		return nil, err
	}

	// Send fee coin to sender
	feeCoin := sdk.NewCoin(params.FeeDenom, amount)

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(feeCoin))
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawUnbondedResponse{}, nil
}
