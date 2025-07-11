package keeper_test

import (
	"testing"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/da/keeper"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestDeleteRejectedDataOvertime(t *testing.T) {
	k, _, _, ctx := setupMsgServer(t)

	// Set removal period
	removalPeriod := 24 * time.Hour

	// Prepare test data
	currentTime := time.Now()
	oldTime := currentTime.Add(-removalPeriod * 2) // 2x removal period ago
	newTime := currentTime.Add(-removalPeriod / 2) // Half removal period ago

	// Set old rejected data
	oldRejectedData := types.PublishedData{
		MetadataUri:       "ipfs://old_rejected",
		Status:            types.Status_STATUS_REJECTED,
		Timestamp:         oldTime,
		ShardDoubleHashes: [][]byte{[]byte("hash")},
	}
	require.NoError(t, k.SetPublishedData(ctx, oldRejectedData))

	// Set new rejected data
	newRejectedData := types.PublishedData{
		MetadataUri:       "ipfs://new_rejected",
		Status:            types.Status_STATUS_REJECTED,
		Timestamp:         newTime,
		ShardDoubleHashes: [][]byte{[]byte("hash")},
	}
	require.NoError(t, k.SetPublishedData(ctx, newRejectedData))

	// Execute test
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := k.DeleteRejectedDataOvertime(sdkCtx, removalPeriod)
	require.NoError(t, err)

	// Verify old data has been deleted
	_, found, err := k.GetPublishedData(ctx, oldRejectedData.MetadataUri)
	require.NoError(t, err)
	require.False(t, found)

	// Verify new data is still present
	_, found, err = k.GetPublishedData(ctx, newRejectedData.MetadataUri)
	require.NoError(t, err)
	require.True(t, found)
}

func TestDeleteVerifiedDataOvertime(t *testing.T) {
	k, _, _, ctx := setupMsgServer(t)

	// Set removal period
	removalPeriod := 24 * time.Hour

	// Prepare test data
	currentTime := time.Now()
	oldTime := currentTime.Add(-removalPeriod * 2)
	newTime := currentTime.Add(-removalPeriod / 2)

	// Set old verified data
	oldVerifiedData := types.PublishedData{
		MetadataUri:       "ipfs://old_verified",
		Status:            types.Status_STATUS_VERIFIED,
		Timestamp:         oldTime,
		ShardDoubleHashes: [][]byte{[]byte("hash")},
	}
	require.NoError(t, k.SetPublishedData(ctx, oldVerifiedData))

	// Set new verified data
	newVerifiedData := types.PublishedData{
		MetadataUri:       "ipfs://new_verified",
		Status:            types.Status_STATUS_VERIFIED,
		Timestamp:         newTime,
		ShardDoubleHashes: [][]byte{[]byte("hash")},
	}
	require.NoError(t, k.SetPublishedData(ctx, newVerifiedData))

	// Execute test
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := k.DeleteVerifiedDataOvertime(sdkCtx, removalPeriod)
	require.NoError(t, err)

	// Verify old data has been deleted
	_, found, err := k.GetPublishedData(ctx, oldVerifiedData.MetadataUri)
	require.NoError(t, err)
	require.False(t, found)

	// Verify new data is still present
	_, found, err = k.GetPublishedData(ctx, newVerifiedData.MetadataUri)
	require.NoError(t, err)
	require.True(t, found)
}

func TestChangeToChallengingFromChallengePeriod(t *testing.T) {
	k, mocks, _, ctx := setupMsgServer(t)

	// Set challenge threshold
	threshold := "0.5" // 50% threshold

	// Prepare test data
	currentTime := time.Now()
	shardDoubleHashes := [][]byte{[]byte("hash1"), []byte("hash2")}

	// Set data in challenge period with sufficient invalidity
	dataWithEnoughInvalidity := types.PublishedData{
		MetadataUri:       "ipfs://enough_invalidity",
		Status:            types.Status_STATUS_CHALLENGE_PERIOD,
		Timestamp:         currentTime,
		ShardDoubleHashes: shardDoubleHashes,
	}
	require.NoError(t, k.SetPublishedData(ctx, dataWithEnoughInvalidity))

	// Add invalidity (above threshold)
	sender := sdk.AccAddress("sender")
	invalidity1 := types.Invalidity{
		MetadataUri: dataWithEnoughInvalidity.MetadataUri,
		Sender:      sender.String(),
		Indices:     []int64{0},
	}
	require.NoError(t, k.SetInvalidity(ctx, invalidity1))

	// Set data in challenge period with insufficient invalidity
	dataWithInsufficientInvalidity := types.PublishedData{
		MetadataUri:       "ipfs://insufficient_invalidity",
		Status:            types.Status_STATUS_CHALLENGE_PERIOD,
		Timestamp:         currentTime,
		ShardDoubleHashes: shardDoubleHashes,
	}
	require.NoError(t, k.SetPublishedData(ctx, dataWithInsufficientInvalidity))

	// Set up mock for staking keeper
	mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(new(mockIterator), nil).AnyTimes()
	mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), gomock.Any()).Return(nil, stakingtypes.ErrNoValidatorFound).AnyTimes()

	// Execute test
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := k.ChangeToChallengingFromChallengePeriod(sdkCtx, threshold)
	require.NoError(t, err)

	// Verify status has changed for data with sufficient invalidity
	data, found, err := k.GetPublishedData(ctx, dataWithEnoughInvalidity.MetadataUri)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, types.Status_STATUS_CHALLENGING, data.Status)

	// Verify status has not changed for data with insufficient invalidity
	data, found, err = k.GetPublishedData(ctx, dataWithInsufficientInvalidity.MetadataUri)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, types.Status_STATUS_CHALLENGE_PERIOD, data.Status)
}

func TestChangeToVerifiedFromProofPeriod(t *testing.T) {
	k, mocks, _, ctx := setupMsgServer(t)

	// Set challenge period
	challengePeriod := 24 * time.Hour

	// Prepare test data
	currentTime := time.Now()
	oldTime := currentTime.Add(-challengePeriod * 2)
	newTime := currentTime.Add(-challengePeriod / 2)

	publisher := sdk.AccAddress("publisher")
	collateral := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))

	// Set expired challenge period data
	expiredData := types.PublishedData{
		MetadataUri:           "ipfs://expired_challenge",
		Status:                types.Status_STATUS_CHALLENGE_PERIOD,
		Timestamp:             oldTime,
		Publisher:             publisher.String(),
		PublishDataCollateral: collateral,
		ShardDoubleHashes:     [][]byte{[]byte("hash")},
	}
	require.NoError(t, k.SetPublishedData(ctx, expiredData))

	// Set valid challenge period data
	validData := types.PublishedData{
		MetadataUri:           "ipfs://valid_challenge",
		Status:                types.Status_STATUS_CHALLENGE_PERIOD,
		Timestamp:             newTime,
		Publisher:             publisher.String(),
		PublishDataCollateral: collateral,
		ShardDoubleHashes:     [][]byte{[]byte("hash")},
	}
	require.NoError(t, k.SetPublishedData(ctx, validData))

	// Set up mock for coin transfer
	mocks.BankKeeper.EXPECT().
		SendCoinsFromModuleToAccount(ctx, types.ModuleName, publisher, collateral).
		Return(nil)

	// Execute test
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := k.ChangeToVerifiedFromChallengePeriod(sdkCtx, challengePeriod)
	require.NoError(t, err)

	// Verify status has changed for expired data
	data, found, err := k.GetPublishedData(ctx, expiredData.MetadataUri)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, types.Status_STATUS_VERIFIED, data.Status)

	// Verify status has not changed for valid data
	data, found, err = k.GetPublishedData(ctx, validData.MetadataUri)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, types.Status_STATUS_CHALLENGE_PERIOD, data.Status)
}

func TestTallyValidityProofs(t *testing.T) {
	// Mock validator
	pubKey1 := ed25519.GenPrivKey().PubKey()
	anyPk1, err := codectypes.NewAnyWithValue(pubKey1)
	require.NoError(t, err)
	valAddr1 := sdk.ValAddress(pubKey1.Address())
	mockVal1 := &stakingtypes.Validator{
		OperatorAddress: valAddr1.String(),
		ConsensusPubkey: anyPk1,
		Status:          stakingtypes.Bonded,
	}

	pubKey2 := ed25519.GenPrivKey().PubKey()
	anyPk2, err := codectypes.NewAnyWithValue(pubKey2)
	require.NoError(t, err)
	valAddr2 := sdk.ValAddress(pubKey2.Address())
	mockVal2 := &stakingtypes.Validator{
		OperatorAddress: valAddr2.String(),
		ConsensusPubkey: anyPk2,
		Status:          stakingtypes.Bonded,
	}

	publisher := sdk.AccAddress("publisher")
	challenger := sdk.AccAddress("challenger")
	collateral := sdk.NewCoins(sdk.NewInt64Coin("stake", 1000))
	challengerCollateral := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))

	testCases := []struct {
		name                   string
		metadataUri            string
		setup                  func(k keeper.Keeper, mocks DaMocks, ctx sdk.Context)
		replicationFactor      string
		duration               time.Duration
		expectedStatus         types.Status
		expectedFaultCount     map[string]uint64 // map[valAddr]count
		expectPublisherRefund  bool
		expectChallengerReward bool
	}{
		{
			name:        "success - data verified, one validator at fault",
			metadataUri: "ipfs://verified",
			setup: func(k keeper.Keeper, mocks DaMocks, ctx sdk.Context) {
				// Data in challenging state
				data := types.PublishedData{
					MetadataUri:                "ipfs://verified",
					Status:                     types.Status_STATUS_CHALLENGING,
					Timestamp:                  time.Now().Add(-2 * time.Hour),
					ShardDoubleHashes:          [][]byte{[]byte("hash1"), []byte("hash2")},
					ParityShardCount:           0,
					ChallengingValidators:      []string{valAddr1.String(), valAddr2.String()},
					Publisher:                  publisher.String(),
					PublishDataCollateral:      collateral,
					SubmitInvalidityCollateral: challengerCollateral,
				}
				require.NoError(t, k.SetPublishedData(ctx, data))

				// Proofs submitted by only one of the two assigned validators
				proof := types.Proof{
					MetadataUri: data.MetadataUri,
					Sender:      sdk.AccAddress(valAddr1).String(),
					Indices:     []int64{0, 1},
				}
				require.NoError(t, k.SetProof(ctx, proof))

				// Mock staking keeper calls
				mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(&mockIterator{valAddrs: []sdk.ValAddress{valAddr1, valAddr2}}, nil).AnyTimes()
				mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr1).Return(mockVal1, nil).AnyTimes()
				mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr2).Return(mockVal2, nil).AnyTimes()

				// Mock bank keeper for publisher refund
				mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, publisher, collateral).Return(nil)
			},
			replicationFactor:     "1.0",
			duration:              1 * time.Hour,
			expectedStatus:        types.Status_STATUS_VERIFIED,
			expectedFaultCount:    map[string]uint64{valAddr2.String(): 1}, // val2 did not submit proof
			expectPublisherRefund: true,
		},
		{
			name:        "success - data rejected, not enough proofs",
			metadataUri: "ipfs://rejected",
			setup: func(k keeper.Keeper, mocks DaMocks, ctx sdk.Context) {
				data := types.PublishedData{
					MetadataUri:                "ipfs://rejected",
					Status:                     types.Status_STATUS_CHALLENGING,
					Timestamp:                  time.Now().Add(-2 * time.Hour),
					ShardDoubleHashes:          [][]byte{[]byte("hash1"), []byte("hash2")},
					ParityShardCount:           0,
					ChallengingValidators:      []string{valAddr1.String(), valAddr2.String()},
					Publisher:                  publisher.String(),
					PublishDataCollateral:      collateral,
					SubmitInvalidityCollateral: challengerCollateral,
				}
				require.NoError(t, k.SetPublishedData(ctx, data))

				// No proofs submitted

				// One challenger
				invalidity := types.Invalidity{MetadataUri: data.MetadataUri, Sender: challenger.String()}
				require.NoError(t, k.SetInvalidity(ctx, invalidity))

				// Mock staking keeper calls
				mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(&mockIterator{valAddrs: []sdk.ValAddress{valAddr1, valAddr2}}, nil).AnyTimes()
				mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr1).Return(mockVal1, nil).AnyTimes()
				mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr2).Return(mockVal2, nil).AnyTimes()

				// Mock bank keeper for challenger reward
				// Reward = challenger's own collateral + publisher's collateral
				totalReward := challengerCollateral.Add(collateral...)
				mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, challenger, totalReward).Return(nil)
			},
			replicationFactor:      "2.0", // Requires more proofs than available
			duration:               1 * time.Hour,
			expectedStatus:         types.Status_STATUS_REJECTED,
			expectedFaultCount:     map[string]uint64{}, // No one is at fault if data is rejected
			expectChallengerReward: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k, mocks, _, ctx := setupMsgServer(t)
			sdkCtx := sdk.UnwrapSDKContext(ctx)

			tc.setup(k, mocks, sdkCtx)

			err := k.TallyValidityProofs(sdkCtx, tc.duration, tc.replicationFactor)
			require.NoError(t, err)

			// Verify data status
			data, found, err := k.GetPublishedData(ctx, tc.metadataUri)
			require.NoError(t, err)
			require.True(t, found)
			require.Equal(t, tc.expectedStatus, data.Status)

			// Verify fault counters
			for valAddr, expectedCount := range tc.expectedFaultCount {
				addr, _ := sdk.ValAddressFromBech32(valAddr)
				count, err := k.GetFaultCounter(ctx, addr)
				require.NoError(t, err)
				require.Equal(t, expectedCount, count)
			}
		})
	}
}
