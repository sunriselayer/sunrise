package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestMsgPublishData(t *testing.T) {
	k, mocks, ms, ctx := setupMsgServer(t)
	params := types.DefaultParams()
	// set default collateral
	params.PublishDataCollateral = sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	params.SubmitInvalidityCollateral = sdk.NewCoins(sdk.NewInt64Coin("stake", 50))
	require.NoError(t, k.Params.Set(ctx, params))

	sender := sdk.AccAddress("sender")
	validMetadataUri := "ipfs://metadata1"
	shardDoubleHashes := [][]byte{[]byte("hash1"), []byte("hash2"), []byte("hash3")}

	// set existing data
	existingMetadataUri := "ipfs://existing"
	err := k.SetPublishedData(ctx, types.PublishedData{
		MetadataUri:                existingMetadataUri,
		ParityShardCount:           1,
		ShardDoubleHashes:          shardDoubleHashes,
		Timestamp:                  time.Now(),
		Status:                     types.Status_STATUS_CHALLENGE_PERIOD,
		Publisher:                  sender.String(),
		PublishDataCollateral:      params.PublishDataCollateral,
		SubmitInvalidityCollateral: params.SubmitInvalidityCollateral,
		PublishedTimestamp:         time.Now(),
	})
	require.NoError(t, err)

	testCases := []struct {
		name      string
		input     *types.MsgPublishData
		mockSetup func()
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid sender address",
			input: &types.MsgPublishData{
				Sender:            "invalid_address",
				MetadataUri:       validMetadataUri,
				ParityShardCount:  1,
				ShardDoubleHashes: shardDoubleHashes,
			},
			expErr:    true,
			expErrMsg: "invalid sender address",
		},
		{
			name: "parity shard count is greater than total",
			input: &types.MsgPublishData{
				Sender:            sender.String(),
				MetadataUri:       validMetadataUri,
				ParityShardCount:  uint64(len(shardDoubleHashes)),
				ShardDoubleHashes: shardDoubleHashes,
			},
			expErr:    true,
			expErrMsg: "parity shard count is greater than total",
		},
		{
			name: "data already exist",
			input: &types.MsgPublishData{
				Sender:            sender.String(),
				MetadataUri:       existingMetadataUri,
				ParityShardCount:  1,
				ShardDoubleHashes: shardDoubleHashes,
			},
			expErr:    true,
			expErrMsg: "data already exist",
		},
		{
			name: "normal case",
			input: &types.MsgPublishData{
				Sender:            sender.String(),
				MetadataUri:       validMetadataUri,
				ParityShardCount:  1,
				ShardDoubleHashes: shardDoubleHashes,
				DataSourceInfo:    "test data source",
			},
			mockSetup: func() {
				mocks.BankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expErr: false,
		},
		{
			name: "failed to send coins",
			input: &types.MsgPublishData{
				Sender:            sender.String(),
				MetadataUri:       "ipfs://another_metadata",
				ParityShardCount:  1,
				ShardDoubleHashes: shardDoubleHashes,
			},
			mockSetup: func() {
				mocks.BankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdkerrors.ErrInsufficientFunds)
			},
			expErr:    true,
			expErrMsg: "insufficient funds",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockSetup != nil {
				tc.mockSetup()
			}

			_, err := ms.PublishData(ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)

				// in normal case, published data is saved in storage
				publishedData, found, err := k.GetPublishedData(ctx, tc.input.MetadataUri)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, tc.input.MetadataUri, publishedData.MetadataUri)
				require.Equal(t, tc.input.ParityShardCount, publishedData.ParityShardCount)
				require.Equal(t, tc.input.ShardDoubleHashes, publishedData.ShardDoubleHashes)
				require.Equal(t, tc.input.Sender, publishedData.Publisher)
				require.Equal(t, tc.input.DataSourceInfo, publishedData.DataSourceInfo)
				require.Equal(t, types.Status_STATUS_CHALLENGE_PERIOD, publishedData.Status)
				require.Equal(t, params.PublishDataCollateral, publishedData.PublishDataCollateral)
				require.Equal(t, params.SubmitInvalidityCollateral, publishedData.SubmitInvalidityCollateral)
				require.False(t, publishedData.Timestamp.IsZero())
				require.False(t, publishedData.PublishedTimestamp.IsZero())
			}
		})
	}
}
