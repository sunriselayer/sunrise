package customkeeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	customtypes "github.com/sunriselayer/sunrise/x/bank/types"
)

type Keeper struct {
	bankkeeper.BaseKeeper
	ak              banktypes.AccountKeeper
	BeforeSendHooks []customtypes.BeforeSendHook
}

func NewKeeper(
	bankKeeper bankkeeper.BaseKeeper,
	ak banktypes.AccountKeeper,
	beforeSendHooks []customtypes.BeforeSendHook,
) Keeper {
	return Keeper{
		BaseKeeper:      bankKeeper,
		ak:              ak,
		BeforeSendHooks: beforeSendHooks,
	}
}

// <exora>
// SendCoins transfers amt coins from a sending account to a receiving account.
// An error is returned upon failure.
func (k Keeper) SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	for _, hook := range k.BeforeSendHooks {
		if err := hook(ctx, fromAddr, toAddr, amt); err != nil {
			return err
		}
	}

	return k.BaseKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

// </exora>

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress.
// It will panic if the module account does not exist. An error is returned if
// the recipient address is black-listed or if sending the tokens fails.
func (k Keeper) SendCoinsFromModuleToAccount(
	ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins,
) error {
	senderAddr := k.ak.GetModuleAddress(senderModule)
	if senderAddr == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	if k.BlockedAddr(recipientAddr) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", recipientAddr)
	}

	return k.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another.
// It will panic if either module account does not exist.
func (k Keeper) SendCoinsFromModuleToModule(
	ctx context.Context, senderModule, recipientModule string, amt sdk.Coins,
) error {
	senderAddr := k.ak.GetModuleAddress(senderModule)
	if senderAddr == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return k.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount.
// It will panic if the module account does not exist.
func (k Keeper) SendCoinsFromAccountToModule(
	ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins,
) error {
	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return k.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}
