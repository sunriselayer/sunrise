package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/tokenconverter/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) WithdrawSelfUnbonded(ctx context.Context, msg *types.MsgWithdrawSelfUnbonded) (*types.MsgWithdrawSelfUnbondedResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	// end static validation

	accAddress := sdk.AccAddress(sender)

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	proxyModuleName := types.SelfDelegateProxyAccountModuleName(msg.Sender)
	err = k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		proxyModuleName,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, msg.Amount)),
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.BurnCoins(
		ctx,
		k.accountKeeper.GetModuleAddress(types.ModuleName),
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, msg.Amount)),
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.MintCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.FeeDenom, msg.Amount)),
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		accAddress,
		sdk.NewCoins(sdk.NewCoin(params.FeeDenom, msg.Amount)),
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawSelfUnbondedResponse{}, nil
}
