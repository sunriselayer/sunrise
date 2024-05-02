package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) JoinPool(goCtx context.Context, msg *types.MsgJoinPool) (*types.MsgJoinPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address := sdk.MustAccAddressFromBech32(msg.Sender)

	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("pool id %d doesn't exist", msg.PoolId))
	}

	if msg.BaseToken.Denom != pool.BaseDenom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("base denom %s is invalid", msg.BaseToken.Denom))
	}
	if msg.QuoteToken.Denom != pool.QuoteDenom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("quote denom %s is invalid", msg.QuoteToken.Denom))
	}

	x, y := k.GetPoolBalance(ctx, pool)
	price, err := types.CalculatePrice(x, y, pool)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error calculating price: %s", err.Error()))
	}
	value := types.LpTokenValueInQuoteUnit(x, y, *price)

	newX := x.Add(msg.BaseToken.Amount)
	newY := y.Add(msg.QuoteToken.Amount)
	newPrice, err := types.CalculatePrice(newX, newY, pool)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error calculating price: %s", err.Error()))
	}
	newValue := types.LpTokenValueInQuoteUnit(newX, newY, *newPrice)

	if newValue.LTE(value) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "value is not increased")
	}

	supply := k.GetLpTokenSupply(ctx, pool.Id)
	newSupplyAmount := newValue.Quo(value).MulInt(supply.Amount).RoundInt()
	additionalSupply := sdk.NewCoin(supply.Denom, newSupplyAmount.Sub(supply.Amount))

	if additionalSupply.Amount.LT(msg.MinShareAmount) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "min share amount is not met")
	}

	if err := k.TransferFromAccountToPoolModule(ctx, msg.BaseToken, address, pool.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error transferring base token: %s", err.Error()))
	}
	if err := k.TransferFromAccountToPoolModule(ctx, msg.QuoteToken, address, pool.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error transferring quote token: %s", err.Error()))
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(additionalSupply)); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error minting lp token: %s", err.Error()))
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(additionalSupply)); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error sending lp token: %s", err.Error()))
	}

	return &types.MsgJoinPoolResponse{
		ShareAmount: additionalSupply.Amount,
	}, nil
}

func (k msgServer) ExitPool(goCtx context.Context, msg *types.MsgExitPool) (*types.MsgExitPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("pool id %d doesn't exist", msg.PoolId))
	}

	_ = pool

	return &types.MsgExitPoolResponse{}, nil
}
