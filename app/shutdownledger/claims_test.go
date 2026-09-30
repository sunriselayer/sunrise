package shutdownledger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"
	pooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

func TestBuildClaimsSeparatesSourcesAndSkipsProtocolPots(t *testing.T) {
	user := sdk.AccAddress([]byte("claim-user-address01")).String()
	lockOwner := sdk.AccAddress([]byte("lock-owner-address01")).String()
	lockup := sdk.AccAddress([]byte("lockup-account-addr1")).String()
	valA := sdk.ValAddress([]byte("validator-a-address1")).String()
	valB := sdk.ValAddress([]byte("validator-b-address1")).String()
	valC := sdk.ValAddress([]byte("validator-c-address1")).String()
	operatorA := sdk.AccAddress([]byte("validator-a-address1")).String()
	module := authtypes.NewModuleAddress(shareclasstypes.ModuleName).String()
	saver := shareclasstypes.RewardSaverAddress(valB).String()
	poolAddress := pooltypes.NewPoolAddress(1).String()
	feeAddress := pooltypes.NewPoolFeesAddress(1).String()
	escrow := transfertypes.GetEscrowAddress("transfer", "channel-0").String()
	feeCollector := authtypes.NewModuleAddress(authtypes.FeeCollectorName).String()
	shareDenom := shareclasstypes.NonVotingShareTokenDenom(valB)

	snapshot := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	launch := snapshot.AddDate(-1, 0, 0)
	pool := pooltypes.Pool{
		Id:               1,
		DenomBase:        "ibc/ATOM",
		DenomQuote:       denomUUSDRise,
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

	doc := map[string]any{
		"genesis_time":   launch.Format(time.RFC3339),
		"initial_height": "10",
		"app_state": map[string]any{
			"bank": map[string]any{
				"balances": []any{
					map[string]any{"address": user, "coins": []any{
						map[string]any{"denom": denomUVRise, "amount": "6"},
						map[string]any{"denom": denomUUSDRise, "amount": "8"},
						map[string]any{"denom": usdnIBCDenom, "amount": "4"},
						map[string]any{"denom": "ibc/ATOM", "amount": "9"},
						map[string]any{"denom": shareDenom, "amount": "4"},
					}},
					map[string]any{"address": lockup, "coins": []any{
						map[string]any{"denom": denomURise, "amount": "40"},
					}},
					map[string]any{"address": usdriseWrapper, "coins": []any{
						map[string]any{"denom": usdnIBCDenom, "amount": "100"},
						map[string]any{"denom": denomURise, "amount": "3"},
					}},
					map[string]any{"address": poolAddress, "coins": []any{
						map[string]any{"denom": denomUUSDRise, "amount": "999"},
					}},
					map[string]any{"address": feeAddress, "coins": []any{
						map[string]any{"denom": denomUUSDRise, "amount": "50"},
					}},
					map[string]any{"address": escrow, "coins": []any{
						map[string]any{"denom": "ibc/ATOM", "amount": "70"},
					}},
					map[string]any{"address": feeCollector, "coins": []any{
						map[string]any{"denom": denomURise, "amount": "1000"},
					}},
					map[string]any{"address": saver, "coins": []any{
						map[string]any{"denom": denomURise, "amount": "1000"},
					}},
				},
				"supply": []any{
					map[string]any{"denom": shareDenom, "amount": "4"},
				},
			},
			"liquiditypool": map[string]any{
				"pools": []any{pool},
				"positions": []any{
					map[string]any{
						"id": "1", "address": user, "pool_id": "1",
						"lower_tick": "-100", "upper_tick": "100", "liquidity": liquidity.String(),
					},
				},
				"accumulators": []any{
					map[string]any{
						"name": pooltypes.KeyFeePoolAccumulator(1),
						"accum_value": []any{
							map[string]any{"denom": denomUUSDRise, "amount": "2.000000000000000000"},
						},
						"total_shares": "1.000000000000000000",
					},
				},
				"accumulator_positions": []any{
					map[string]any{
						"name":       "fee",
						"index":      pooltypes.KeyFeePositionAccumulator(1),
						"num_shares": "1.000000000000000000",
						"accum_value_per_share": []any{
							map[string]any{"denom": denomUUSDRise, "amount": "0.000000000000000000"},
						},
						"unclaimed_rewards_total": []any{
							map[string]any{"denom": denomUUSDRise, "amount": "1.500000000000000000"},
						},
					},
				},
			},
			"lockup": map[string]any{
				"lockup_accounts": []any{
					map[string]any{
						"address": lockup, "owner": lockOwner, "id": "7",
						"start_time": snapshot.Unix(), "end_time": snapshot.Unix() + 1000,
						"original_locking": "100", "additional_locking": "0",
						"delegated_free": "0", "delegated_locking": "60",
					},
				},
			},
			"swap": map[string]any{
				"incoming_in_flight_packets": []any{},
				"outgoing_in_flight_packets": []any{},
			},
			"shareclass": map[string]any{
				"unbondings": []any{},
				"reward_multipliers": []any{
					map[string]any{"validator": valB, "denom": denomURise, "reward_multiplier": "2.000000000000000000"},
				},
				"user_last_reward_multipliers": []any{
					map[string]any{"user": user, "validator": valB, "denom": denomURise, "reward_multiplier": "1.500000000000000000"},
				},
			},
			"staking": map[string]any{
				"params": map[string]any{"bond_denom": denomUVRise},
				"validators": []any{
					map[string]any{"operator_address": valA, "tokens": "1000", "delegator_shares": "1000.000000000000000000"},
					map[string]any{"operator_address": valB, "tokens": "40", "delegator_shares": "40.000000000000000000"},
					map[string]any{"operator_address": valC, "tokens": "60", "delegator_shares": "60.000000000000000000"},
				},
				"delegations": []any{
					map[string]any{"delegator_address": user, "validator_address": valA, "shares": "100.000000000000000000"},
					map[string]any{"delegator_address": module, "validator_address": valB, "shares": "40.000000000000000000"},
					map[string]any{"delegator_address": lockup, "validator_address": valC, "shares": "60.000000000000000000"},
				},
				"unbonding_delegations": []any{},
			},
			"distribution": map[string]any{
				"validator_historical_rewards": []any{
					map[string]any{
						"validator_address": valA, "period": "0",
						"rewards": map[string]any{"cumulative_reward_ratio": []any{
							map[string]any{"denom": denomURise, "amount": "0.000000000000000000"},
						}},
					},
				},
				"validator_current_rewards": []any{
					map[string]any{
						"validator_address": valA,
						"rewards": map[string]any{
							"period": "1",
							"rewards": []any{
								map[string]any{"denom": denomURise, "amount": "100.000000000000000000"},
							},
						},
					},
				},
				"delegator_starting_infos": []any{
					map[string]any{
						"delegator_address": user, "validator_address": valA,
						"starting_info": map[string]any{
							"previous_period": "0", "stake": "100.000000000000000000", "height": "1",
						},
					},
					map[string]any{
						"delegator_address": module, "validator_address": valB,
						"starting_info": map[string]any{
							"previous_period": "0", "stake": "40.000000000000000000", "height": "10",
						},
					},
					map[string]any{
						"delegator_address": lockup, "validator_address": valC,
						"starting_info": map[string]any{
							"previous_period": "0", "stake": "60.000000000000000000", "height": "10",
						},
					},
				},
				"validator_slash_events": []any{},
				"validator_accumulated_commissions": []any{
					map[string]any{
						"validator_address": valA,
						"accumulated": map[string]any{"commission": []any{
							map[string]any{"denom": denomURise, "amount": "2.500000000000000000"},
						}},
					},
				},
			},
			"ibc": map[string]any{
				"channel_genesis": map[string]any{
					"channels": []any{
						map[string]any{"port_id": "transfer", "channel_id": "channel-0"},
					},
				},
			},
		},
	}

	dir := t.TempDir()
	input := filepath.Join(dir, "pre-upgrade-state.json")
	encoded, err := json.Marshal(doc)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(input, encoded, 0o600))
	output := filepath.Join(dir, "ledger")
	require.NoError(t, Build(input, output, snapshot))

	var file claimsFile
	readJSON(t, filepath.Join(output, ClaimsFile), &file)
	require.Equal(t, snapshot.Unix(), file.SnapshotUnix)
	require.Equal(t, int64(10), file.SnapshotHeight)
	require.NotEmpty(t, file.Slippage)

	require.Equal(t, "6", claimAmount(t, file.Claims, user, SourceBank, denomUVRise))
	require.Equal(t, AssetRise, findClaim(t, file.Claims, user, SourceBank, denomUVRise).Asset)
	require.Equal(t, PayoutEdge, findClaim(t, file.Claims, user, SourceBank, denomUVRise).Payout)
	require.Equal(t, "8", claimAmount(t, file.Claims, user, SourceBank, denomUUSDRise))
	require.Equal(t, AssetUSDrise, findClaim(t, file.Claims, user, SourceBank, denomUUSDRise).Asset)
	require.Equal(t, PayoutEdge, findClaim(t, file.Claims, user, SourceBank, denomUUSDRise).Payout)
	require.Equal(t, "4", claimAmount(t, file.Claims, user, SourceBank, usdnIBCDenom))
	require.Equal(t, AssetUSDrise, findClaim(t, file.Claims, user, SourceBank, usdnIBCDenom).Asset)
	require.Equal(t, "9", claimAmount(t, file.Claims, user, SourceBank, "ibc/ATOM"))
	require.Equal(t, PayoutCosmos, findClaim(t, file.Claims, user, SourceBank, "ibc/ATOM").Payout)

	require.Equal(t, "100", claimAmount(t, file.Claims, user, SourceDelegation, valA))
	require.Equal(t, "10", claimAmount(t, file.Claims, user, SourceStakingReward, valA+":"+denomURise))
	require.Equal(t, "2", claimAmount(t, file.Claims, operatorA, SourceStakingReward, "commission:"+valA+":"+denomURise))
	require.Equal(t, "40", claimAmount(t, file.Claims, user, SourceShareclass, valB))
	require.Equal(t, "2", claimAmount(t, file.Claims, user, SourceStakingReward, "shareclass:"+valB+":"+denomURise))

	require.Equal(t, wantBase.String(), claimAmount(t, file.Claims, user, SourcePosition, "1:base"))
	require.Equal(t, wantQuote.String(), claimAmount(t, file.Claims, user, SourcePosition, "1:quote"))
	require.Equal(t, AssetUSDrise, findClaim(t, file.Claims, user, SourcePosition, "1:quote").Asset)
	require.Equal(t, "3", claimAmount(t, file.Claims, user, SourcePositionReward, "1:"+denomUUSDRise))

	lockedBank := findClaim(t, file.Claims, lockOwner, SourceLockup, "7:bank:"+denomURise+":locked")
	require.Equal(t, "40", lockedBank.Amount)
	require.Equal(t, snapshot.Unix()+1000, lockedBank.ClaimableAt)
	lockedStake := findClaim(t, file.Claims, lockOwner, SourceLockup, "7:delegation:"+valC+":locked")
	require.Equal(t, "60", lockedStake.Amount)
	require.Equal(t, snapshot.Unix()+1000, lockedStake.ClaimableAt)
	require.Equal(t, "3", claimAmount(t, file.Claims, usdriseWrapper, SourceBank, denomURise))

	for _, claim := range file.Claims {
		require.NotEqual(t, poolAddress, claim.Owner)
		require.NotEqual(t, feeAddress, claim.Owner)
		require.NotEqual(t, escrow, claim.Owner)
		require.NotEqual(t, feeCollector, claim.Owner)
		require.NotEqual(t, saver, claim.Owner)
		require.NotEqual(t, module, claim.Owner)
		require.NotEqual(t, lockup, claim.Owner)
		if claim.Owner == usdriseWrapper {
			require.NotEqual(t, AssetUSDrise, claim.Asset)
		}
	}
}

func TestDelegationRewardsIncludeSlashPeriods(t *testing.T) {
	delegator := sdk.AccAddress([]byte("slash-delegator-addr")).String()
	validator := sdk.ValAddress([]byte("slash-validator-addr")).String()
	state, err := prepareRewards(distributionExport{
		Historical: []struct {
			ValidatorAddress string     `json:"validator_address"`
			Period           flexUint64 `json:"period"`
			Rewards          struct {
				Cumulative sdk.DecCoins `json:"cumulative_reward_ratio"`
			} `json:"rewards"`
		}{
			{ValidatorAddress: validator, Period: 0, Rewards: struct {
				Cumulative sdk.DecCoins `json:"cumulative_reward_ratio"`
			}{Cumulative: sdk.DecCoins{sdk.NewDecCoinFromDec(denomURise, math.LegacyZeroDec())}}},
			{ValidatorAddress: validator, Period: 1, Rewards: struct {
				Cumulative sdk.DecCoins `json:"cumulative_reward_ratio"`
			}{Cumulative: sdk.DecCoins{sdk.NewDecCoinFromDec(denomURise, math.LegacyMustNewDecFromStr("0.2"))}}},
		},
		Current: []struct {
			ValidatorAddress string `json:"validator_address"`
			Rewards          struct {
				Rewards sdk.DecCoins `json:"rewards"`
				Period  flexUint64   `json:"period"`
			} `json:"rewards"`
		}{
			{ValidatorAddress: validator, Rewards: struct {
				Rewards sdk.DecCoins `json:"rewards"`
				Period  flexUint64   `json:"period"`
			}{Rewards: sdk.DecCoins{}, Period: 2}},
		},
		Starting: []struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			StartingInfo     struct {
				PreviousPeriod flexUint64     `json:"previous_period"`
				Stake          math.LegacyDec `json:"stake"`
				Height         flexUint64     `json:"height"`
			} `json:"starting_info"`
		}{
			{DelegatorAddress: delegator, ValidatorAddress: validator, StartingInfo: struct {
				PreviousPeriod flexUint64     `json:"previous_period"`
				Stake          math.LegacyDec `json:"stake"`
				Height         flexUint64     `json:"height"`
			}{PreviousPeriod: 0, Stake: math.LegacyNewDec(1000), Height: 1}},
		},
		Slashes: []struct {
			ValidatorAddress string     `json:"validator_address"`
			Height           flexUint64 `json:"height"`
			Event            struct {
				Period   flexUint64     `json:"validator_period"`
				Fraction math.LegacyDec `json:"fraction"`
			} `json:"validator_slash_event"`
		}{
			{ValidatorAddress: validator, Height: 5, Event: struct {
				Period   flexUint64     `json:"validator_period"`
				Fraction math.LegacyDec `json:"fraction"`
			}{Period: 1, Fraction: math.LegacyMustNewDecFromStr("0.1")}},
		},
	}, stakingExport{
		Validators: []struct {
			OperatorAddress string         `json:"operator_address"`
			Tokens          math.Int       `json:"tokens"`
			DelegatorShares math.LegacyDec `json:"delegator_shares"`
		}{
			{OperatorAddress: validator, Tokens: math.NewInt(900), DelegatorShares: math.LegacyNewDec(1000)},
		},
	}, 10)
	require.NoError(t, err)

	coins, err := state.delegationRewards(delegator, validator, math.LegacyNewDec(1000))
	require.NoError(t, err)
	require.Equal(t, "200"+denomURise, coins.String())
}

func TestStartingInfoHeightMatchesSDKGenesisJSON(t *testing.T) {
	genesis := disttypes.DefaultGenesisState()
	genesis.DelegatorStartingInfos = []disttypes.DelegatorStartingInfoRecord{{
		DelegatorAddress: sdk.AccAddress([]byte("height-delegator-adr")).String(),
		ValidatorAddress: sdk.ValAddress([]byte("height-validator-adr")).String(),
		StartingInfo: disttypes.DelegatorStartingInfo{
			PreviousPeriod: 3,
			Stake:          math.LegacyNewDec(100),
			Height:         42,
		},
	}}
	raw, err := moduletestutil.MakeTestEncodingConfig().Codec.MarshalJSON(genesis)
	require.NoError(t, err)

	var distribution distributionExport
	require.NoError(t, json.Unmarshal(raw, &distribution))
	require.Len(t, distribution.Starting, 1)
	require.Equal(t, flexUint64(3), distribution.Starting[0].StartingInfo.PreviousPeriod)
	require.Equal(t, flexUint64(42), distribution.Starting[0].StartingInfo.Height)
}

func TestBuildClaimsVestsLockupAtSnapshotTimeAndSkipsNativeDenoms(t *testing.T) {
	holder := sdk.AccAddress([]byte("native-holder-addr01")).String()
	lockOwner := sdk.AccAddress([]byte("vest-owner-address01")).String()
	lockup := sdk.AccAddress([]byte("vest-lockup-address1")).String()
	factoryDenom := "factory/" + holder + "/meme"
	snapshot := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	launch := snapshot.AddDate(-1, 0, 0)

	doc := map[string]any{
		"genesis_time":   launch.Format(time.RFC3339),
		"initial_height": "20",
		"app_state": map[string]any{
			"bank": map[string]any{
				"balances": []any{
					map[string]any{"address": holder, "coins": []any{
						map[string]any{"denom": factoryDenom, "amount": "5"},
						map[string]any{"denom": "ibc/ATOM", "amount": "7"},
					}},
					map[string]any{"address": lockup, "coins": []any{
						map[string]any{"denom": denomURise, "amount": "80"},
						map[string]any{"denom": denomUVRise, "amount": "30"},
					}},
				},
				"supply": []any{},
			},
			"lockup": map[string]any{
				"lockup_accounts": []any{
					map[string]any{
						"address": lockup, "owner": lockOwner, "id": "3",
						"start_time": snapshot.Unix() - 500, "end_time": snapshot.Unix() + 500,
						"original_locking": "100", "additional_locking": "0",
						"delegated_free": "0", "delegated_locking": "0",
					},
				},
			},
		},
	}

	dir := t.TempDir()
	input := filepath.Join(dir, "pre-upgrade-state.json")
	encoded, err := json.Marshal(doc)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(input, encoded, 0o600))
	require.Error(t, Build(input, filepath.Join(dir, "no-snapshot-time"), time.Time{}))
	output := filepath.Join(dir, "ledger")
	require.NoError(t, Build(input, output, snapshot))

	var file claimsFile
	readJSON(t, filepath.Join(output, ClaimsFile), &file)
	require.Equal(t, snapshot.Unix(), file.SnapshotUnix)

	// Half of the lockup has vested at the snapshot, so 50 of the urise balance stays locked.
	locked := findClaim(t, file.Claims, lockOwner, SourceLockup, "3:bank:"+denomURise+":locked")
	require.Equal(t, "50", locked.Amount)
	require.Equal(t, snapshot.Unix()+500, locked.ClaimableAt)
	unlocked := findClaim(t, file.Claims, lockOwner, SourceLockup, "3:bank:"+denomURise+":unlocked")
	require.Equal(t, "30", unlocked.Amount)
	require.Equal(t, snapshot.Unix(), unlocked.ClaimableAt)
	// The lockup module locks only urise, so uvrise in the account is free.
	vrise := findClaim(t, file.Claims, lockOwner, SourceLockup, "3:bank:"+denomUVRise)
	require.Equal(t, "30", vrise.Amount)
	require.Equal(t, snapshot.Unix(), vrise.ClaimableAt)

	require.Equal(t, PayoutCosmos, findClaim(t, file.Claims, holder, SourceBank, "ibc/ATOM").Payout)
	lockupRows := 0
	for _, claim := range file.Claims {
		require.NotEqual(t, factoryDenom, claim.Asset)
		if claim.Owner == lockOwner {
			lockupRows++
		}
	}
	require.Equal(t, 3, lockupRows)
}

func findClaim(t *testing.T, claims []Claim, owner, source, sourceID string) Claim {
	t.Helper()
	var found []Claim
	for _, claim := range claims {
		if claim.Owner == owner && claim.Source == source && claim.SourceID == sourceID {
			found = append(found, claim)
		}
	}
	require.Len(t, found, 1)
	return found[0]
}

func claimAmount(t *testing.T, claims []Claim, owner, source, sourceID string) string {
	t.Helper()
	return findClaim(t, claims, owner, source, sourceID).Amount
}
