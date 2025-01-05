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

func (k msgServer) SelfUndelegate(ctx context.Context, msg *types.MsgSelfUndelegate) (*types.MsgSelfUndelegateResponse, error) {
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

	// TODO
	var amount math.Int

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	proxyModuleName := types.SelfDelegateProxyAccountModuleName(msg.Creator)
	proxyAddr := k.accountKeeper.GetModuleAddress(proxyModuleName)

	stakingKeeper, ok := k.stakingKeeper.(*stakingkeeper.Keeper)
	if !ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidType, "invalid staking keeper")
	}
	_, err = stakingkeeper.NewMsgServerImpl(stakingKeeper).Undelegate(ctx, &stakingtypes.MsgUndelegate{
		DelegatorAddress: proxyAddr.String(),
		ValidatorAddress: validator.GetOperator(),
		Amount:           sdk.NewCoin(params.BondDenom, amount),
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgSelfDelegateResponse{}, nil
}
