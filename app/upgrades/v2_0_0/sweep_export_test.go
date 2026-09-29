package v2_0_0

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

// TestSweepIBCOnExportedState replays SweepIBC against a sunrised export.
// Set SHUTDOWN_EXPORT to the export path. The test checks that every non-escrow
// ibc balance is sent once, and that native and escrow balances are not sent.
func TestSweepIBCOnExportedState(t *testing.T) {
	path := os.Getenv("SHUTDOWN_EXPORT")
	if path == "" {
		t.Skip("SHUTDOWN_EXPORT is not set")
	}

	raw, err := os.ReadFile(path)
	require.NoError(t, err)

	var doc exportedState
	require.NoError(t, json.Unmarshal(raw, &doc))

	recovery, err := sdk.AccAddressFromBech32(RecoveryAddress)
	require.NoError(t, err)

	channels := make([]channeltypes.IdentifiedChannel, 0)
	escrow := map[string]struct{}{}
	for _, channel := range doc.AppState.IBC.ChannelGenesis.Channels {
		if channel.PortID != transfertypes.PortID {
			continue
		}
		channels = append(channels, channeltypes.IdentifiedChannel{
			PortId:    channel.PortID,
			ChannelId: channel.ChannelID,
		})
		escrow[string(transfertypes.GetEscrowAddress(channel.PortID, channel.ChannelID))] = struct{}{}
	}
	require.NotEmpty(t, channels)

	bank := &memoryBank{}
	total := map[string]math.Int{}
	escrowTotal := map[string]math.Int{}
	recoveryBefore := map[string]math.Int{}
	for _, balance := range doc.AppState.Bank.Balances {
		address, err := sdk.AccAddressFromBech32(balance.Address)
		require.NoError(t, err)
		_, isEscrow := escrow[string(address)]
		isRecovery := address.Equals(recovery)
		for _, coin := range balance.Coins {
			amount, ok := math.NewIntFromString(coin.Amount)
			require.True(t, ok, "amount %s for %s", coin.Amount, coin.Denom)
			bank.rows = append(bank.rows, balanceRow{address: address, coin: sdk.NewCoin(coin.Denom, amount)})
			if !strings.HasPrefix(coin.Denom, IBCDenomPrefix) {
				continue
			}
			total[coin.Denom] = addAmount(total[coin.Denom], amount)
			if isEscrow {
				escrowTotal[coin.Denom] = addAmount(escrowTotal[coin.Denom], amount)
			}
			if isRecovery {
				recoveryBefore[coin.Denom] = addAmount(recoveryBefore[coin.Denom], amount)
			}
		}
	}

	ctx := sdk.Context{}.WithLogger(log.NewNopLogger())
	require.NoError(t, SweepIBC(ctx, bank, channelList{channels: channels}))

	swept := map[string]math.Int{}
	seenFrom := map[string]struct{}{}
	for _, sent := range bank.sent {
		require.Equal(t, RecoveryAddress, sent.to)
		require.NotContains(t, escrow, string(mustAccAddress(t, sent.from)))
		require.NotEqual(t, RecoveryAddress, sent.from)
		_, dup := seenFrom[sent.from]
		require.False(t, dup, "account swept twice: %s", sent.from)
		seenFrom[sent.from] = struct{}{}
		for _, coin := range sent.coins {
			require.True(t, strings.HasPrefix(coin.Denom, IBCDenomPrefix))
			swept[coin.Denom] = addAmount(swept[coin.Denom], coin.Amount)
		}
	}

	for denom, supply := range total {
		got := addAmount(swept[denom], addAmount(escrowTotal[denom], recoveryBefore[denom]))
		require.Truef(t, supply.Equal(got), "denom %s supply %s swept+escrow+recovery %s", denom, supply, got)
	}
	t.Logf("accounts=%d swept_accounts=%d channels=%d", len(doc.AppState.Bank.Balances), len(bank.sent), len(channels))
	for denom, amount := range swept {
		t.Logf("swept %s %s escrow %s recovery_before %s", denom, amount, escrowTotal[denom], recoveryBefore[denom])
	}
}

func addAmount(current, amount math.Int) math.Int {
	if amount.IsNil() {
		if current.IsNil() {
			return math.ZeroInt()
		}
		return current
	}
	if current.IsNil() {
		return amount
	}
	return current.Add(amount)
}

func mustAccAddress(t *testing.T, address string) sdk.AccAddress {
	t.Helper()
	parsed, err := sdk.AccAddressFromBech32(address)
	require.NoError(t, err)
	return parsed
}

type exportedState struct {
	AppState struct {
		Bank struct {
			Balances []struct {
				Address string `json:"address"`
				Coins   []struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"coins"`
			} `json:"balances"`
		} `json:"bank"`
		IBC struct {
			ChannelGenesis struct {
				Channels []struct {
					PortID    string `json:"port_id"`
					ChannelID string `json:"channel_id"`
				} `json:"channels"`
			} `json:"channel_genesis"`
		} `json:"ibc"`
	} `json:"app_state"`
}
