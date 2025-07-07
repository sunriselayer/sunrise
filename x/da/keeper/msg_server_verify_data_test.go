package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestMsgVerifyData(t *testing.T) {
	k, mocks, msgServer, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// --- Test Data Setup ---
	challengePeriod := 24 * time.Hour
	removalPeriod := 48 * time.Hour
	publisher := sdk.AccAddress("publisher")
	challenger := sdk.AccAddress("challenger")
	collateral := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))

	// 1. Data that will expire and be verified
	expiredChallengeData := types.PublishedData{
		MetadataUri:           "ipfs://expired_challenge",
		Status:                types.Status_STATUS_CHALLENGE_PERIOD,
		Timestamp:             sdkCtx.BlockTime().Add(-challengePeriod * 2),
		Publisher:             publisher.String(),
		PublishDataCollateral: collateral,
		ShardDoubleHashes:     [][]byte{[]byte("hash")},
	}
	require.NoError(t, k.SetPublishedData(ctx, expiredChallengeData))

	// 2. Data that will be removed (rejected)
	oldRejectedData := types.PublishedData{
		MetadataUri:       "ipfs://old_rejected",
		Status:            types.Status_STATUS_REJECTED,
		Timestamp:         sdkCtx.BlockTime().Add(-removalPeriod * 2),
		ShardDoubleHashes: [][]byte{[]byte("hash")},
	}
	require.NoError(t, k.SetPublishedData(ctx, oldRejectedData))

	// 3. Data that will become challenging
	challengingData := types.PublishedData{
		MetadataUri:       "ipfs://to_be_challenging",
		Status:            types.Status_STATUS_CHALLENGE_PERIOD,
		Timestamp:         sdkCtx.BlockTime(),
		ShardDoubleHashes: [][]byte{[]byte("hash1"), []byte("hash2")},
	}
	require.NoError(t, k.SetPublishedData(ctx, challengingData))
	invalidity := types.Invalidity{
		MetadataUri: challengingData.MetadataUri,
		Sender:      challenger.String(),
		Indices:     []int64{0, 1}, // Assuming 2 shards, this is > 50%
	}
	require.NoError(t, k.SetInvalidity(ctx, invalidity))

	// --- Mock Expectations ---
	mocks.BankKeeper.EXPECT().
		SendCoinsFromModuleToAccount(sdkCtx, types.ModuleName, publisher, collateral).
		Return(nil).
		AnyTimes()

	// --- Set Params ---
	params := types.DefaultParams()
	params.ChallengePeriod = challengePeriod
	params.RejectedRemovalPeriod = removalPeriod
	params.VerifiedRemovalPeriod = removalPeriod
	params.ChallengeThreshold = "0.5"
	require.NoError(t, k.Params.Set(ctx, params))

	// --- Execute MsgVerifyData ---
	_, err := msgServer.VerifyData(ctx, &types.MsgVerifyData{Sender: publisher.String()})
	require.NoError(t, err)

	// --- Verifications ---
	// 1. Check if expired data is now verified
	data, found, err := k.GetPublishedData(ctx, expiredChallengeData.MetadataUri)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, types.Status_STATUS_VERIFIED, data.Status)

	// 2. Check if old rejected data has been deleted
	_, found, err = k.GetPublishedData(ctx, oldRejectedData.MetadataUri)
	require.NoError(t, err)
	require.False(t, found)

	// 3. Check if data became challenging
	data, found, err = k.GetPublishedData(ctx, challengingData.MetadataUri)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, types.Status_STATUS_CHALLENGING, data.Status)
}
