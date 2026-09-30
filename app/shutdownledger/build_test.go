package shutdownledger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	pooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestMain(m *testing.M) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("sunrise", "sunrisepub")
	config.SetBech32PrefixForValidator("sunrisevaloper", "sunrisevaloperpub")
	config.SetBech32PrefixForConsensusNode("sunrisevalcons", "sunrisevalconspub")
	os.Exit(m.Run())
}

func TestBuildSeparatesBankBalancesFromPositions(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner-address-012345")).String()
	poolAddress := sdk.AccAddress([]byte("pool-address-0012345")).String()
	lockupAddress := sdk.AccAddress([]byte("lockup-address-01234")).String()
	lockupOwner := sdk.AccAddress([]byte("lock-owner-address-01")).String()

	pool := pooltypes.Pool{
		Id:               1,
		DenomBase:        "ibc/BASE",
		DenomQuote:       "urise",
		CurrentTick:      0,
		CurrentSqrtPrice: "1",
		TickParams: pooltypes.TickParams{
			PriceRatio: "1.0001",
			BaseOffset: "0",
		},
	}
	liquidity := math.LegacyNewDec(1_000_000)
	wantBase, wantQuote, err := withdrawalAmounts(pool, -100, 100, liquidity)
	require.NoError(t, err)
	require.True(t, wantBase.IsPositive())
	require.True(t, wantQuote.IsPositive())

	doc := map[string]any{
		"app_state": map[string]any{
			"bank": map[string]any{
				"balances": []any{
					map[string]any{"address": owner, "coins": []any{map[string]any{"denom": "urise", "amount": "15"}}},
					map[string]any{"address": poolAddress, "coins": []any{map[string]any{"denom": "ibc/BASE", "amount": "999"}}},
					map[string]any{"address": lockupAddress, "coins": []any{map[string]any{"denom": "urise", "amount": "40"}}},
				},
				"supply": []any{},
			},
			"liquiditypool": map[string]any{
				"pools": []any{pool},
				"positions": []any{
					map[string]any{
						"id": uint64(1), "address": owner, "pool_id": uint64(1),
						"lower_tick": int64(-100), "upper_tick": int64(100), "liquidity": liquidity.String(),
					},
				},
				"accumulator_positions": []any{
					map[string]any{
						"name": "fee", "index": "fee_position_accumulator/|1", "num_shares": "1",
						"unclaimed_rewards_total": []any{map[string]any{"denom": "urise", "amount": "1.500000000000000000"}},
					},
				},
			},
			"lockup": map[string]any{
				"lockup_accounts": []any{
					map[string]any{"address": lockupAddress, "owner": lockupOwner, "id": uint64(2)},
				},
			},
			"swap": map[string]any{
				"incoming_in_flight_packets": []any{},
				"outgoing_in_flight_packets": []any{
					map[string]any{
						"index":             map[string]any{"port_id": "transfer", "channel_id": "channel-0", "sequence": uint64(9631)},
						"ack_waiting_index": map[string]any{"port_id": "transfer", "channel_id": "channel-1", "sequence": uint64(10)},
						"retries_remaining": int32(2),
					},
				},
			},
			"shareclass": map[string]any{"unbondings": []any{}},
			"staking": map[string]any{
				"params":                map[string]any{"bond_denom": "uvrise"},
				"validators":            []any{},
				"delegations":           []any{},
				"unbonding_delegations": []any{},
			},
		},
	}
	dir := t.TempDir()
	input := filepath.Join(dir, "pre-upgrade-state.json")
	encoded, err := json.Marshal(doc)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(input, encoded, 0o600))

	output := filepath.Join(dir, "ledger")
	require.NoError(t, Build(input, output, time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)))

	var bank struct {
		Balances []struct {
			Address string `json:"address"`
			Coins   []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"coins"`
		} `json:"balances"`
	}
	readJSON(t, filepath.Join(output, BankFile), &bank)
	require.Len(t, bank.Balances, 3)
	for _, balance := range bank.Balances {
		if balance.Address == poolAddress {
			require.Equal(t, "999", balance.Coins[0].Amount)
			require.Equal(t, "ibc/BASE", balance.Coins[0].Denom)
		}
	}

	var positions []positionRecord
	readJSON(t, filepath.Join(output, PositionsFile), &positions)
	require.Len(t, positions, 1)
	require.Equal(t, owner, positions[0].Owner)
	require.Equal(t, wantBase.String(), positions[0].AmountBase)
	require.Equal(t, wantQuote.String(), positions[0].AmountQuote)
	require.Equal(t, "ibc/BASE", positions[0].DenomBase)
	require.Equal(t, "urise", positions[0].DenomQuote)
	require.Equal(t, "1.500000000000000000urise", positions[0].UnclaimedRewards.String())

	var lockups []lockupRecord
	readJSON(t, filepath.Join(output, LockupFile), &lockups)
	require.Len(t, lockups, 1)
	require.Equal(t, lockupOwner, lockups[0].Owner)
	require.Equal(t, lockupAddress, lockups[0].LockupAddress)
	require.Equal(t, "40urise", lockups[0].Balances.String())

	var inFlight inFlightFile
	readJSON(t, filepath.Join(output, InFlightFile), &inFlight)
	require.Empty(t, inFlight.Incoming)
	require.Len(t, inFlight.Outgoing, 1)
	require.Equal(t, uint64(9631), inFlight.Outgoing[0].Sequence)
	require.Equal(t, int32(2), inFlight.Outgoing[0].RetriesRemaining)
}

func readJSON(t *testing.T, path string, dest any) {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, dest))
}
