package keeper

import (
	"context"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// YbtbrandBeforeSendHook returns a BeforeSendHook for ybtbrand transfers
func YbtbrandBeforeSendHook(k Keeper) func(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	return func(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
		// Process each coin in the transfer
		for _, coin := range amount {
			// Check if this is a ybtbrand token
			if !strings.HasPrefix(coin.Denom, "ybtbrand/") {
				continue
			}

			// Extract creator from denom
			parts := strings.Split(coin.Denom, "/")
			if len(parts) != 2 {
				continue
			}
			creator := parts[1]

			// Get token info to check if it exists
			_, found := k.GetToken(ctx, creator)
			if !found {
				continue
			}

			// Get all yield denoms for this token
			// We need to iterate through all possible yield indexes
			yieldDenoms := k.getActiveYieldDenoms(ctx, creator)

			for _, yieldDenom := range yieldDenoms {
				// Get sender's last yield reward index
				fromYieldIndex, _ := k.GetUserLastYieldIndex(ctx, creator, from.String(), yieldDenom)

				// Get receiver's current balance and last yield reward index
				toBalance := k.bankKeeper.GetBalance(ctx, to, coin.Denom)
				toYieldIndex, _ := k.GetUserLastYieldIndex(ctx, creator, to.String(), yieldDenom)

				// Calculate weighted average of receiver's new yield reward index
				// newYieldIndex = (receiverBalance * receiverYieldIndex + transferAmount * senderYieldIndex) / (receiverBalance + transferAmount)
				// This formula ensures that:
				// 1. Rewards are fairly distributed based on when tokens were acquired
				// 2. Negative indexes are properly handled (when sender's index < receiver's index)
				
				// Convert to Dec for calculation
				receiverBalanceDec := math.LegacyNewDecFromInt(toBalance.Amount)
				transferAmountDec := math.LegacyNewDecFromInt(coin.Amount)
				totalAmountDec := receiverBalanceDec.Add(transferAmountDec)

				var newYieldIndex math.LegacyDec
				if totalAmountDec.IsZero() {
					// Edge case: if total is zero, keep receiver's yield index
					newYieldIndex = toYieldIndex
				} else {
					// Calculate weighted average
					weightedReceiverYieldIndex := receiverBalanceDec.Mul(toYieldIndex)
					weightedSenderYieldIndex := transferAmountDec.Mul(fromYieldIndex)
					newYieldIndex = weightedReceiverYieldIndex.Add(weightedSenderYieldIndex).Quo(totalAmountDec)
				}

				// Update receiver's last yield reward index
				if err := k.SetUserLastYieldIndex(ctx, creator, to.String(), yieldDenom, newYieldIndex); err != nil {
					return err
				}
			}
		}

		return nil
	}
}

// getActiveYieldDenoms returns all yield denoms that have been used for this token
func (k Keeper) getActiveYieldDenoms(ctx context.Context, creator string) []string {
	// Iterate through all yield indexes to find active denoms
	denoms := make([]string, 0)
	
	// Use the YieldIndex collection to find all denoms
	iter, err := k.YieldIndex.Iterate(ctx, collections.NewPrefixedPairRange[string, string](creator))
	if err != nil {
		return denoms
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		kv, err := iter.KeyValue()
		if err != nil {
			continue
		}
		// The key is a pair of (creator, denom)
		key := kv.Key
		if key.K1() == creator {
			denoms = append(denoms, key.K2())
		}
	}

	return denoms
}