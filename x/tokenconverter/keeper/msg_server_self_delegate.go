package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/tokenconverter/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "cosmossdk.io/x/staking/types"
)

func (k msgServer) SelfDelegate(ctx context.Context, msg *types.MsgSelfDelegate) (*types.MsgSelfDelegateResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
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

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	proxyModuleName := types.SelfDelegateProxyAccountModuleName(msg.Sender)
	proxyAddr := k.accountKeeper.GetModuleAddress(proxyModuleName)

	bonded, err := k.stakingKeeper.GetDelegatorBonded(ctx, accAddress)
	if err != nil {
		return nil, err
	}

	if bonded.Add(msg.Amount).GT(params.SelfDelegationCap) {
		return nil, errorsmod.Wrapf(types.ErrExceedSelfDelegationCap, "%s + %s > %s", bonded.String(), msg.Amount.String(), params.SelfDelegationCap.String())
	}

	// To support vesting account, use delegate
	err = k.bankKeeper.DelegateCoinsFromAccountToModule(
		ctx,
		accAddress,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.FeeDenom, msg.Amount)),
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.BurnCoins(
		ctx,
		k.accountKeeper.GetModuleAddress(types.ModuleName),
		sdk.NewCoins(sdk.NewCoin(params.FeeDenom, msg.Amount)),
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.MintCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, msg.Amount)),
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		proxyModuleName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, msg.Amount)),
	)
	if err != nil {
		return nil, err
	}

	_, err = k.MsgRouterService.Invoke(ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: proxyAddr.String(),
		ValidatorAddress: validator.GetOperator(),
		Amount:           sdk.NewCoin(params.BondDenom, msg.Amount),
	})
	if err != nil {
		return nil, err
	}

	err = k.distributionKeeper.SetWithdrawAddr(ctx, proxyAddr, accAddress)
	if err != nil {
		return nil, err
	}

	return &types.MsgSelfDelegateResponse{}, nil
}
