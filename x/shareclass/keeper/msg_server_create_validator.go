package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/shareclass/types"

	banktypes "cosmossdk.io/x/bank/types"
	stakingtypes "cosmossdk.io/x/staking/types"
)

func (k msgServer) CreateValidator(ctx context.Context, msg *types.MsgCreateValidator) (*types.MsgCreateValidatorResponse, error) {
	address, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	params, err := k.Keeper.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	if !msg.Fee.Equal(params.CreateValidatorFee) {
		return nil, errorsmod.Wrap(types.ErrInvalidCreateValidatorFee, "invalid create validator fee")
	}

	// Burn fee - 1 urise
	burnCoin := sdk.NewCoin(msg.Fee.Denom, msg.Fee.Amount.Sub(math.OneInt()))

	_, err = k.Environment.MsgRouterService.Invoke(ctx, &banktypes.MsgBurn{
		FromAddress: sdk.AccAddress(address).String(),
		Amount:      []*sdk.Coin{&burnCoin},
	})
	if err != nil {
		return nil, err
	}

	// Convert 1 urise to 1 uvrise
	err = k.tokenConverterKeeper.ConvertReverse(ctx, math.OneInt(), sdk.AccAddress(address))
	if err != nil {
		return nil, err
	}

	// Stake 1 uvrise
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}

	_, err = k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgCreateValidator{
		Description:       msg.Description,
		Commission:        msg.Commission,
		MinSelfDelegation: msg.MinSelfDelegation,
		ValidatorAddress:  msg.ValidatorAddress,
		Pubkey:            msg.Pubkey,
		Value:             sdk.NewCoin(bondDenom, math.OneInt()),
	})

	if err != nil {
		return nil, err
	}

	return &types.MsgCreateValidatorResponse{}, nil
}
