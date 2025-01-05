package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/tokenconverter/types"

	errorsmod "cosmossdk.io/errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) WithdrawSelfUnbonded(ctx context.Context, msg *types.MsgWithdrawSelfUnbonded) (*types.MsgWithdrawSelfUnbonded, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	// end static validation

	accAddress := sdk.AccAddress(sender)
	// TODO
	var amount math.Int

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	proxyModuleName := types.SelfDelegateProxyAccountModuleName(msg.Creator)
	err = k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		proxyModuleName,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, amount)),
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.BurnCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, amount)),
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.MintCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.FeeDenom, amount)),
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		accAddress,
		sdk.NewCoins(sdk.NewCoin(params.FeeDenom, amount)),
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgSelfDelegateResponse{}, nil
}
