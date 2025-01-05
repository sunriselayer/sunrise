package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/tokenconverter/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stakingkeeper "cosmossdk.io/x/staking/keeper"
	stakingtypes "cosmossdk.io/x/staking/types"
)

func (k msgServer) SelfCancelUnbonding(ctx context.Context, msg *types.MsgSelfCancelUnbonding) (*types.MsgSelfCancelUnbondingResponse, error) {
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

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	proxyModuleName := types.SelfDelegateProxyAccountModuleName(msg.Sender)
	proxyAddr := k.accountKeeper.GetModuleAddress(proxyModuleName)

	stakingKeeper, ok := k.stakingKeeper.(*stakingkeeper.Keeper)
	if !ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidType, "invalid staking keeper")
	}
	_, err = stakingkeeper.NewMsgServerImpl(stakingKeeper).CancelUnbondingDelegation(ctx, &stakingtypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: proxyAddr.String(),
		ValidatorAddress: validator.GetOperator(),
		Amount:           sdk.NewCoin(params.BondDenom, msg.Amount),
		CreationHeight:   msg.CreationHeight,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgSelfCancelUnbondingResponse{}, nil
}
