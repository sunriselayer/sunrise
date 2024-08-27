package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestPublishedDataStore(t *testing.T) {
	k, _, _, _, ctx := keepertest.DaKeeper(t)
	sender1 := sdk.AccAddress("sender1")

	data := k.GetPublishedData(ctx, "ipfs://metadata1")
	require.Equal(t, data.MetadataUri, "")

	now := time.Now()
	dataArr := []types.PublishedData{
		{
			MetadataUri:        "ipfs://metadata1",
			ParityShardCount:   0,
			ShardDoubleHashes:  [][]byte{[]byte("data1")},
			Timestamp:          now,
			Status:             "msg_server",
			Publisher:          sender1.String(),
			Challenger:         "",
			Collateral:         sdk.Coins{},
			ChallengeTimestamp: time.Time{},
			DataSourceInfo:     "",
		},
		{
			MetadataUri:        "ipfs://metadata2",
			ParityShardCount:   0,
			ShardDoubleHashes:  [][]byte{[]byte("data1")},
			Timestamp:          now,
			Status:             "vote_extension",
			Publisher:          sender1.String(),
			Challenger:         "",
			Collateral:         sdk.Coins{},
			ChallengeTimestamp: time.Time{},
			DataSourceInfo:     "",
		},
		{
			MetadataUri:        "ipfs://metadata3",
			ParityShardCount:   0,
			ShardDoubleHashes:  [][]byte{[]byte("data1")},
			Timestamp:          now.Add(-time.Second),
			Status:             "challenge_for_fraud",
			Publisher:          sender1.String(),
			Challenger:         "",
			Collateral:         sdk.Coins{},
			ChallengeTimestamp: time.Time{},
			DataSourceInfo:     "",
		},
		{
			MetadataUri:        "ipfs://metadata4",
			ParityShardCount:   0,
			ShardDoubleHashes:  [][]byte{[]byte("data1")},
			Timestamp:          now.Add(-time.Second),
			Status:             "rejected",
			Publisher:          sender1.String(),
			Challenger:         "",
			Collateral:         sdk.Coins{},
			ChallengeTimestamp: time.Time{},
			DataSourceInfo:     "",
		},
		{
			MetadataUri:        "ipfs://metadata5",
			ParityShardCount:   0,
			ShardDoubleHashes:  [][]byte{[]byte("data1")},
			Timestamp:          now,
			Status:             "verified",
			Publisher:          sender1.String(),
			Challenger:         "",
			Collateral:         sdk.Coins{},
			ChallengeTimestamp: time.Time{},
			DataSourceInfo:     "",
		},
	}

	for _, data := range dataArr {
		err := k.SetPublishedData(ctx, data)
		require.NoError(t, err)
	}

	for _, data := range dataArr {
		stored := k.GetPublishedData(ctx, data.MetadataUri)
		require.Equal(t, stored.MetadataUri, data.MetadataUri)
		require.Equal(t, stored.ParityShardCount, data.ParityShardCount)
		require.Equal(t, stored.Timestamp.Unix(), data.Timestamp.Unix())
		require.Equal(t, stored.Status, data.Status)
	}

	dataArr, err := k.GetUnverifiedDataBeforeTime(ctx, uint64(now.Unix()))
	require.NoError(t, err)
	require.Len(t, dataArr, 1)

	dataArr, err = k.GetUnverifiedDataBeforeTime(ctx, uint64(now.Add(-time.Second).Unix()))
	require.NoError(t, err)
	require.Len(t, dataArr, 0)

	dataArr = k.GetAllPublishedData(ctx)
	require.Len(t, dataArr, 5)

	k.DeletePublishedData(ctx, dataArr[2])
	data = k.GetPublishedData(ctx, dataArr[2].MetadataUri)
	require.Equal(t, data.MetadataUri, "")

	dataArr, err = k.GetUnverifiedDataBeforeTime(ctx, uint64(now.Unix()))
	require.NoError(t, err)
	require.Len(t, dataArr, 0)

	dataArr = k.GetAllPublishedData(ctx)
	require.Len(t, dataArr, 4)
}
