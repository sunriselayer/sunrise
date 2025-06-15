package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func (k msgServer) AddYields(ctx context.Context, msg *types.MsgAddYields) (*types.MsgAddYieldsResponse, error) {
	// Validate admin address
	adminAddr, err := k.addressCodec.StringToBytes(msg.Admin)
	if err != nil {
		return nil, errors.Wrap(err, "invalid admin address")
	}

	// Validate token creator address
	if _, err := k.addressCodec.StringToBytes(msg.TokenCreator); err != nil {
		return nil, errors.Wrap(err, "invalid token creator address")
	}

	// Get token
	token, found := k.Keeper.GetToken(ctx, msg.TokenCreator)
	if !found {
		return nil, types.ErrTokenNotFound
	}

	// Check if sender is admin
	if token.Admin != msg.Admin {
		return nil, types.ErrUnauthorized
	}

	// Validate amount
	if len(msg.Amount) == 0 {
		return nil, errors.Wrap(types.ErrInvalidRequest, "amount cannot be empty")
	}
	coins := sdk.NewCoins(msg.Amount...)
	if !coins.IsValid() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "invalid amount")
	}

	// Get brand token supply for yield index calculation
	brandDenom := GetTokenDenom(msg.TokenCreator)
	supply := k.bankKeeper.GetSupply(ctx, brandDenom)

	// Process each yield token
	for _, coin := range msg.Amount {
		if !coin.IsPositive() {
			return nil, errors.Wrapf(types.ErrInvalidRequest, "yield amount must be positive: %s", coin.String())
		}

		// Transfer yield to yield pool
		yieldPoolAddr := GetYieldPoolAddress(msg.TokenCreator, coin.Denom)
		if err := k.bankKeeper.SendCoins(ctx, sdk.AccAddress(adminAddr), yieldPoolAddr, sdk.NewCoins(coin)); err != nil {
			return nil, err
		}

		// Create yield pool module account if it doesn't exist
		if acc := k.authKeeper.GetAccount(ctx, yieldPoolAddr); acc == nil {
			acc := authtypes.NewModuleAccount(
				authtypes.NewBaseAccountWithAddress(yieldPoolAddr),
				yieldPoolAddr.String(),
			)
			k.authKeeper.SetModuleAccount(ctx, acc)
		}

		// Update yield index if supply is not zero
		if !supply.IsZero() {
			currentIndex, _ := k.Keeper.GetYieldIndex(ctx, msg.TokenCreator, coin.Denom)
			// If index doesn't exist, start from 1.0
			if currentIndex.IsZero() {
				currentIndex = math.LegacyOneDec()
			}

			// Calculate yield per token
			yieldPerToken := math.LegacyNewDecFromInt(coin.Amount).Quo(math.LegacyNewDecFromInt(supply.Amount))
			newIndex := currentIndex.Add(yieldPerToken)

			if err := k.Keeper.SetYieldIndex(ctx, msg.TokenCreator, coin.Denom, newIndex); err != nil {
				return nil, err
			}
		} else {
			// If no supply, just ensure index is initialized at 1.0
			if _, found := k.Keeper.GetYieldIndex(ctx, msg.TokenCreator, coin.Denom); !found {
				if err := k.Keeper.SetYieldIndex(ctx, msg.TokenCreator, coin.Denom, math.LegacyOneDec()); err != nil {
					return nil, err
				}
			}
		}
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddYields,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.TokenCreator),
			sdk.NewAttribute(types.AttributeKeyAdmin, msg.Admin),
			sdk.NewAttribute(types.AttributeKeyAmount, coins.String()),
		),
	})

	return &types.MsgAddYieldsResponse{}, nil
}
