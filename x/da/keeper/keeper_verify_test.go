package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/da/types"
)

type mockIterator struct {
	values [][]byte
	idx    int
}

func (m *mockIterator) Domain() (start, end []byte) {
	return nil, nil
}

func (m *mockIterator) Valid() bool {
	return m.idx < len(m.values)
}

func (m *mockIterator) Next() {
	if m.Valid() {
		m.idx++
	}
}

func (m *mockIterator) Key() (key []byte) {
	return m.values[m.idx]
}

func (m *mockIterator) Value() (value []byte) {
	return m.values[m.idx]
}

func (m *mockIterator) Error() error {
	return nil
}

func (m *mockIterator) Close() error {
	return nil
}

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
