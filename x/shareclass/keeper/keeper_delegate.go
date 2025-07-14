package keeper

import (
	"context"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k Keeper) Delegate(ctx context.Context, sender sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) (share, rewards sdk.Coins, err error) {
	// Validate amount
	transferableDenom, err := k.tokenConverterKeeper.GetTransferableDenom(ctx)
	if err != nil {
		return nil, nil, err
	}
	if amount.Denom != transferableDenom {
		return nil, nil, errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid denom: expected %s, got %s", transferableDenom, amount.Denom)
	}

	// Claim rewards
	rewards, err = k.ClaimRewards(ctx, sender, valAddr)
	if err != nil {
		return nil, nil, err
	}

	// Calculate share before delegate
	shareAmount, err := k.CalculateShareByAmount(ctx, valAddr.String(), amount.Amount)
	if err != nil {
		return nil, nil, err
	}

	// Convert and delegate
	err = k.ConvertAndDelegate(ctx, sender, valAddr, amount.Amount)
	if err != nil {
		return nil, nil, err
	}

	// Mint non transferrable share token
	shareDenom := types.NonVotingShareTokenDenom(valAddr.String())
	k.bankKeeper.SetSendEnabled(ctx, shareDenom, false)
	share = sdk.NewCoins(sdk.NewCoin(shareDenom, shareAmount))

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, share)
	if err != nil {
		return nil, nil, err
	}

	// Send non transferrable share token to sender
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, share)
	if err != nil {
		return nil, nil, err
	}

	return share, rewards, nil
}

func (k Keeper) Undelegate(ctx context.Context, sender sdk.AccAddress, recipient sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) (output sdk.Coin, rewards sdk.Coins, CompletionTime time.Time, err error) {
	// Validate params and amount
	tokenconverterParams, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return sdk.Coin{}, nil, time.Time{}, err
	}
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return sdk.Coin{}, nil, time.Time{}, err
	}
	if tokenconverterParams.NonTransferableDenom != bondDenom {
		return sdk.Coin{}, nil, time.Time{}, errorsmod.Wrapf(types.ErrInvalidTokenConverterParams, "invalid token converter non transferable denom: expected %s, got %s", bondDenom, tokenconverterParams.NonTransferableDenom)
	}
	if amount.Denom != tokenconverterParams.TransferableDenom {
		return sdk.Coin{}, nil, time.Time{}, errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid denom: expected %s, got %s", tokenconverterParams.TransferableDenom, amount.Denom)
	}
	if !amount.IsPositive() {
		return sdk.Coin{}, nil, time.Time{}, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "undelegate amount must be positive")
	}

	// Claim rewards
	rewards, err = k.ClaimRewards(ctx, sender, valAddr)
	if err != nil {
		return sdk.Coin{}, nil, time.Time{}, err
	}

	// Calculate unbonding share
	shareDenom := types.NonVotingShareTokenDenom(valAddr.String())
	unbondingShare, err := k.CalculateShareByAmount(ctx, valAddr.String(), amount.Amount)
	if err != nil {
		return sdk.Coin{}, nil, time.Time{}, err
	}

	// Send non transferrable share token to module
	coins := sdk.NewCoins(sdk.NewCoin(shareDenom, unbondingShare))
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins)
	if err != nil {
		return sdk.Coin{}, nil, time.Time{}, err
	}

	// Burn non transferrable share token
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return sdk.Coin{}, nil, time.Time{}, err
	}

	// Undelegate
	amount = sdk.NewCoin(bondDenom, amount.Amount)

	res, err := k.StakingMsgServer.Undelegate(ctx, &stakingtypes.MsgUndelegate{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
	})
	if err != nil {
		return sdk.Coin{}, nil, time.Time{}, err
	}
	completionTime := res.CompletionTime
	output = res.Amount

	// Append Unstaking state
	_, err = k.AppendUnbonding(ctx, types.Unbonding{
		RecipientAddress: recipient.String(),
		DelegatorAddress: sender.String(),
		ValidatorAddress: valAddr.String(),
		CompletionTime:   completionTime,
		Amount:           output,
	})
	if err != nil {
		return sdk.Coin{}, nil, time.Time{}, err
	}

	return output, rewards, completionTime, nil
}

func (k Keeper) ConvertAndDelegate(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress, amount math.Int) error {
	// Prepare fee and bond coin
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return err
	}
	tokenconverterParams, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return err
	}
	if tokenconverterParams.NonTransferableDenom != bondDenom {
		return errorsmod.Wrapf(types.ErrInvalidTokenConverterParams, "invalid token converter non transferable denom: expected %s, got %s", bondDenom, tokenconverterParams.NonTransferableDenom)
	}
	bondCoin := sdk.NewCoin(bondDenom, amount)
	transferableCoin := sdk.NewCoin(tokenconverterParams.TransferableDenom, amount)

	// Send fee coin to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(transferableCoin))
	if err != nil {
		return err
	}

	// Convert fee denom to bond denom
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	err = k.tokenConverterKeeper.ConvertReverse(ctx, amount, moduleAddr)
	if err != nil {
		return err
	}

	// Stake
	_, err = k.StakingMsgServer.Delegate(ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Amount:           bondCoin,
	})
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetTotalStakedAmount(ctx context.Context, validatorAddr string) (math.Int, error) {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	res, err := k.StakingQueryServer.Delegation(ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: moduleAddr.String(),
		ValidatorAddr: validatorAddr,
	})
	if err != nil {
		return math.Int{}, err
	}

	stakedAmount := res.DelegationResponse.Balance.Amount

	return stakedAmount, nil
}
