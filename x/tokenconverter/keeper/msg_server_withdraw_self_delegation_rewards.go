package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/tokenconverter/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	distributionkeeper "cosmossdk.io/x/distribution/keeper"
	distributiontypes "cosmossdk.io/x/distribution/types"
)

func (k msgServer) WithdrawSelfDelegationRewards(ctx context.Context, msg *types.MsgWithdrawSelfDelegationRewards) (*types.MsgSelfDelegateResponse, error) {
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

	proxyModuleName := types.SelfDelegateProxyAccountModuleName(msg.Sender)
	proxyAddr := k.accountKeeper.GetModuleAddress(proxyModuleName)

	distributionKeeper, ok := k.distributionKeeper.(distributionkeeper.Keeper)
	if !ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidType, "invalid distribution keeper")
	}
	_, err = distributionkeeper.NewMsgServerImpl(distributionKeeper).WithdrawDelegatorReward(ctx, &distributiontypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: proxyAddr.String(),
		ValidatorAddress: validator.GetOperator(),
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgSelfDelegateResponse{}, nil
}
