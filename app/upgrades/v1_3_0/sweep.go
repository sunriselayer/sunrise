package v1_3_0

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
)

// BankKeeper moves IBC vouchers during the upgrade.
type BankKeeper interface {
	IterateAllBalances(ctx context.Context, cb func(address sdk.AccAddress, coin sdk.Coin) bool)
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// ChannelLister lists transfer channels so their escrow accounts can be left in place.
type ChannelLister interface {
	GetAllChannelsWithPortPrefix(ctx sdk.Context, portPrefix string) []channeltypes.IdentifiedChannel
}

// SweepIBC sends ibc/ balances to the recovery account.
// Transfer escrow accounts and the recovery account itself are skipped.
// Native denoms, including urise and uusdrise, stay where they are.
func SweepIBC(ctx sdk.Context, bank BankKeeper, channels ChannelLister) error {
	recovery, err := sdk.AccAddressFromBech32(RecoveryAddress)
	if err != nil {
		return fmt.Errorf("SweepIBC: parse recovery address %s: %w", RecoveryAddress, err)
	}

	skip := escrowAccountKeys(ctx, channels)
	skip[string(recovery)] = struct{}{}

	type transfer struct {
		from  sdk.AccAddress
		coins sdk.Coins
	}
	pending := make([]transfer, 0)
	var current sdk.AccAddress
	var coins sdk.Coins

	flush := func() {
		if len(coins) == 0 || len(current) == 0 {
			coins = nil
			return
		}
		pending = append(pending, transfer{from: current, coins: coins})
		coins = nil
	}

	bank.IterateAllBalances(ctx, func(address sdk.AccAddress, coin sdk.Coin) bool {
		if !current.Equals(address) {
			flush()
			current = address
		}
		if !shouldSweepCoin(address, coin, skip) {
			return false
		}
		coins = coins.Add(coin)
		return false
	})
	flush()

	for _, item := range pending {
		if err := bank.SendCoins(ctx, item.from, recovery, item.coins); err != nil {
			return fmt.Errorf("SweepIBC: SendCoins from %s amount %s to %s: %w", item.from.String(), item.coins.String(), RecoveryAddress, err)
		}
		ctx.Logger().Info("swept IBC vouchers", "from", item.from.String(), "amount", item.coins.String(), "to", RecoveryAddress)
	}
	ctx.Logger().Info("finished IBC sweep", "accounts", len(pending), "recovery", RecoveryAddress)
	return nil
}

func escrowAccountKeys(ctx sdk.Context, channels ChannelLister) map[string]struct{} {
	skip := make(map[string]struct{})
	if channels == nil {
		return skip
	}
	for _, channel := range channels.GetAllChannelsWithPortPrefix(ctx, transfertypes.PortID) {
		if channel.PortId != transfertypes.PortID {
			continue
		}
		escrow := transfertypes.GetEscrowAddress(channel.PortId, channel.ChannelId)
		skip[string(escrow)] = struct{}{}
	}
	return skip
}

func shouldSweepCoin(address sdk.AccAddress, coin sdk.Coin, skip map[string]struct{}) bool {
	if coin.Amount.IsNil() || !coin.Amount.IsPositive() {
		return false
	}
	if !strings.HasPrefix(coin.Denom, IBCDenomPrefix) {
		return false
	}
	if _, found := skip[string(address)]; found {
		return false
	}
	return true
}
