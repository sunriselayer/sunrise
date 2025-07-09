package keeper

import (
	"context"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

// BuyLockupAccount handles MsgBuyLockupAccount.
// It allows a buyer to purchase a listed lockup account.
func (k msgServer) BuyLockupAccount(ctx context.Context, msg *types.MsgBuyLockupAccount) (*types.MsgBuyLockupAccountResponse, error) {
	buyer, err := k.accountKeeper.AddressCodec().StringToBytes(msg.Buyer)
	if err != nil {
		return nil, err
	}
	seller, err := k.accountKeeper.AddressCodec().StringToBytes(msg.Seller)
	if err != nil {
		return nil, err
	}

	if msg.Buyer == msg.Seller {
		return nil, types.ErrInvalidSale
	}

	// Get the listing price
	price, err := k.Listings.Get(ctx, collections.Join(seller, msg.LockupAccountId))
	if err != nil {
		return nil, types.ErrNotListed
	}

	// Transfer funds from buyer to seller
	if err := k.bankKeeper.SendCoins(ctx, buyer, seller, sdk.NewCoins(price)); err != nil {
		return nil, err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	lockupAccount, err := k.LockupAccounts.Get(ctx, collections.Join(seller, msg.LockupAccountId))
	if err != nil {
		return nil, err
	}

	// Calculate and pay incentive
	incentiveRate, err := math.LegacyNewDecFromStr(params.MarketIncentiveRate)
	if err != nil {
		return nil, err
	}
	incentiveAmountDec := math.LegacyNewDecFromInt(lockupAccount.OriginalLocking).Mul(incentiveRate)
	incentiveAmount := incentiveAmountDec.TruncateInt()

	maxIncentive := math.NewIntFromUint64(params.MaxMarketIncentiveAmount)
	if incentiveAmount.GT(maxIncentive) {
		incentiveAmount = maxIncentive
	}

	incentiveCoin := sdk.NewCoin(feeDenom, incentiveAmount)
	if incentiveCoin.IsPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.IncentivePoolName, seller, sdk.NewCoins(incentiveCoin)); err != nil {
			return nil, err
		}
	}

	// Transfer ownership
	if err := k.LockupAccounts.Remove(ctx, collections.Join(seller, msg.LockupAccountId)); err != nil {
		return nil, err
	}

	newLockupId, _, err := k.GetAndIncrementNextLockupAccountID(ctx, buyer)
	if err != nil {
		return nil, err
	}

	lockupAccount.Owner = msg.Buyer
	lockupAccount.Id = newLockupId
	lockupAccount.EndTime = time.Unix(lockupAccount.EndTime, 0).Add(params.MarketLockupDurationExtension).Unix()

	if err := k.LockupAccounts.Set(ctx, collections.Join(buyer, newLockupId), lockupAccount); err != nil {
		return nil, err
	}

	// Remove listing
	if err := k.Listings.Remove(ctx, collections.Join(seller, msg.LockupAccountId)); err != nil {
		return nil, err
	}

	return &types.MsgBuyLockupAccountResponse{}, nil
}
