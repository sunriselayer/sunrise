package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/tokenconverter/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stakingkeeper "cosmossdk.io/x/staking/keeper"
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
	// TODO
	var amount math.Int

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		accAddress,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.FeeDenom, amount)),
	)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.BurnCoins(
		ctx,
		types.ModuleName,
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

	proxyModuleName := types.SelfDelegateProxyAccountModuleName(msg.Sender)
	err = k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		proxyModuleName,
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, amount)),
	)
	if err != nil {
		return nil, err
	}

	proxyAddr := k.accountKeeper.GetModuleAddress(proxyModuleName)

	stakingKeeper, ok := k.stakingKeeper.(*stakingkeeper.Keeper)
	if !ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidType, "invalid staking keeper")
	}
	_, err = stakingkeeper.NewMsgServerImpl(stakingKeeper).Delegate(ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: proxyAddr.String(),
		ValidatorAddress: validator.GetOperator(),
		Amount:           sdk.NewCoin(params.BondDenom, amount),
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
