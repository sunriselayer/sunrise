package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/shareclass/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k msgServer) CreateValidator(ctx context.Context, msg *types.MsgCreateValidator) (*types.MsgCreateValidatorResponse, error) {
	address, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	// Validate amount
	powerReduction := k.stakingKeeper.PowerReduction(ctx)

	if !msg.Amount.Amount.Equal(powerReduction) {
		return nil, errorsmod.Wrapf(types.ErrInvalidCreateValidatorAmount, "invalid create validator amount: expected %s (power reduction in staking module), got %s", powerReduction, msg.Amount.Amount)
	}

	tokenconverterParams, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}
	if tokenconverterParams.NonTransferableDenom != bondDenom {
		return nil, errorsmod.Wrapf(types.ErrInvalidTokenConverterParams, "invalid token converter non transferable denom: expected %s, got %s", bondDenom, tokenconverterParams.NonTransferableDenom)
	}
	if msg.Amount.Denom != tokenconverterParams.TransferableDenom {
		return nil, errorsmod.Wrapf(types.ErrInvalidCreateValidatorAmount, "invalid denom: expected %s, got %s", tokenconverterParams.TransferableDenom, msg.Amount.Denom)
	}

	// Consume gas
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	sdk.UnwrapSDKContext(ctx).GasMeter().ConsumeGas(params.CreateValidatorGas, "create validator with fee denom")

	// Convert amount from fee denom to bond denom
	err = k.tokenConverterKeeper.ConvertReverse(ctx, msg.Amount.Amount, sdk.AccAddress(address))
	if err != nil {
		return nil, err
	}

	// MsgCreateValidator in cosmos-sdk reads CachedValue, so create it here.
	var pk cryptotypes.PubKey
	err = k.cdc.UnpackAny(msg.Pubkey, &pk)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to unpack public key from Any")
	}

	pkAny, err := codectypes.NewAnyWithValue(pk)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to pack public key into Any")
	}

	_, err = k.StakingMsgServer.CreateValidator(ctx, &stakingtypes.MsgCreateValidator{
		Description:       msg.Description,
		Commission:        msg.Commission,
		MinSelfDelegation: msg.MinSelfDelegation,
		ValidatorAddress:  msg.ValidatorAddress,
		Pubkey:            pkAny,
		Value:             sdk.NewCoin(bondDenom, msg.Amount.Amount),
	})

	if err != nil {
		return nil, err
	}

	return &types.MsgCreateValidatorResponse{}, nil
}
