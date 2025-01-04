package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/tokenconverter/types"

	errorsmod "cosmossdk.io/errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "cosmossdk.io/x/staking/types"
)

func (k msgServer) SelfDelegate(ctx context.Context, msg *types.MsgSelfDelegate) (*types.MsgSelfDelegateResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	// end static validation

	valAddress := sdk.ValAddress(sender)
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddress)
	if err != nil {
		return nil, err
	}

	accAddress := sdk.AccAddress(sender)
	var amount math.Int

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		accAddress,
		types.SelfDelegateProxyAccountModuleName(msg.Creator),
		sdk.NewCoins(sdk.NewCoin(params.FeeDenom, amount)),
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.MintCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, amount)),
	)
	if err != nil {
		return nil, err
	}

	proxyModuleName := types.SelfDelegateProxyAccountModuleName(msg.Creator)
	err = k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		proxyModuleName,
		sdk.NewCoins(sdk.NewCoin(params.FeeDenom, amount)),
	)
	if err != nil {
		return nil, err
	}

	proxyAddr := k.accountKeeper.GetModuleAddress(proxyModuleName)

	_, err = k.stakingKeeper.Delegate(ctx, proxyAddr, amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	return &types.MsgSelfDelegateResponse{}, nil
}
