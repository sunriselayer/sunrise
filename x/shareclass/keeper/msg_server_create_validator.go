package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
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

	// Validate amount
	powerReduction := k.stakingKeeper.PowerReduction(ctx)

	if !msg.Amount.Amount.Equal(powerReduction) {
		return nil, errorsmod.Wrap(types.ErrInvalidCreateValidatorAmount, "create validator amount must be equal to power reduction in staking module")
	}

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Denom != feeDenom {
		return nil, errorsmod.Wrap(types.ErrInvalidCreateValidatorAmount, "create validator amount denom must be equal to fee denom")
	}

	// Validate fee
	params, err := k.Keeper.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	if !msg.Fee.Equal(params.CreateValidatorFee) {
		return nil, errorsmod.Wrap(types.ErrInvalidCreateValidatorFee, "invalid create validator fee")
	}

	// Burn fee
	_, err = k.Environment.MsgRouterService.Invoke(ctx, &banktypes.MsgBurn{
		FromAddress: sdk.AccAddress(address).String(),
		Amount:      []*sdk.Coin{&msg.Fee},
	})
	if err != nil {
		return nil, err
	}

	// Convert amount from fee denom to bond denom
	err = k.tokenConverterKeeper.ConvertReverse(ctx, msg.Amount.Amount, sdk.AccAddress(address))
	if err != nil {
		return nil, err
	}

	// Stake
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
		Value:             sdk.NewCoin(bondDenom, msg.Amount.Amount),
	})

	if err != nil {
		return nil, err
	}

	return &types.MsgCreateValidatorResponse{}, nil
}
