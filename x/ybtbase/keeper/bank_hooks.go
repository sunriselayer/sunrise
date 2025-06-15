package keeper

import (
	"context"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// YbtbaseBeforeSendHook returns a BeforeSendHook for ybtbase transfers
func YbtbaseBeforeSendHook(k Keeper) func(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	return func(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
		// Process each coin in the transfer
		for _, coin := range amount {
			// Check if this is a ybtbase token
			if !strings.HasPrefix(coin.Denom, "ybtbase/") {
				continue
			}

			// Extract creator from denom
			parts := strings.Split(coin.Denom, "/")
			if len(parts) != 2 {
				continue
			}
			creator := parts[1]

			// Get token info to check if it exists
			token, found := k.GetToken(ctx, creator)
			if !found {
				continue
			}

			// Get sender's last reward index
			fromRewardIndex := k.GetUserLastRewardIndex(ctx, creator, from.String())

			// Get receiver's current balance and last reward index
			toBalance := k.bankKeeper.GetBalance(ctx, to, coin.Denom)
			toRewardIndex := k.GetUserLastRewardIndex(ctx, creator, to.String())

			// Calculate weighted average of receiver's new reward index
			// newRewardIndex = (receiverBalance * receiverRewardIndex + transferAmount * senderRewardIndex) / (receiverBalance + transferAmount)
			// This formula ensures that:
			// 1. Rewards are fairly distributed based on when tokens were acquired
			// 2. Negative indexes are properly handled (when sender's index < receiver's index)

			// Convert to Dec for calculation
			receiverBalanceDec := math.LegacyNewDecFromInt(toBalance.Amount)
			transferAmountDec := math.LegacyNewDecFromInt(coin.Amount)
			totalAmountDec := receiverBalanceDec.Add(transferAmountDec)

			var newRewardIndex math.LegacyDec
			if totalAmountDec.IsZero() {
				// Edge case: if total is zero, keep receiver's reward index
				newRewardIndex = toRewardIndex
			} else {
				// Calculate weighted average
				weightedReceiverRewardIndex := receiverBalanceDec.Mul(toRewardIndex)
				weightedSenderRewardIndex := transferAmountDec.Mul(fromRewardIndex)
				newRewardIndex = weightedReceiverRewardIndex.Add(weightedSenderRewardIndex).Quo(totalAmountDec)
			}

			// Update receiver's last reward index
			if err := k.SetUserLastRewardIndex(ctx, creator, to.String(), newRewardIndex); err != nil {
				return err
			}

			// For permissioned tokens, check if transfer should grant permission
			if token.Permissioned {
				// If sender has permission and receiver doesn't, consider granting
				// This is a policy decision - we'll leave permission unchanged
				// Admin must explicitly grant permissions
			}
		}

		return nil
	}
}
