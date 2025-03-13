package keeper_test

import (
	"math/big"
	"testing"
	"time"

	native_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestMsgSubmitInvalidity(t *testing.T) {
	k, mocks, ms, ctx := setupMsgServer(t)
	params := types.DefaultParams()
	// set challenge period to 1 day
	params.ChallengePeriod = time.Hour * 24
	require.NoError(t, k.Params.Set(ctx, params))

	sender := sdk.AccAddress("sender")

	// set valid metadata uri to published data
	validMetadataUri := "ipfs://metadata1"
	currentTime := time.Now()

	preImage1 := big.NewInt(111)
	m := native_mimc.NewMiMC()
	m.Write(preImage1.Bytes())
	hash := m.Sum(nil)

	// set valid metadata uri to published data
	err := k.SetPublishedData(ctx, types.PublishedData{
		MetadataUri:                validMetadataUri,
		ParityShardCount:           0,
		ShardDoubleHashes:          [][]byte{hash},
		Timestamp:                  currentTime,
		Status:                     types.Status_STATUS_CHALLENGE_PERIOD,
		Publisher:                  "publisher",
		PublishDataCollateral:      sdk.Coins{},
		SubmitInvalidityCollateral: sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
		PublishedTimestamp:         currentTime,
	})
	require.NoError(t, err)

	// set expired data
	expiredMetadataUri := "ipfs://expired"
	expiredTime := currentTime.Add(-params.ChallengePeriod * 2)
	err = k.SetPublishedData(ctx, types.PublishedData{
		MetadataUri:                expiredMetadataUri,
		ParityShardCount:           0,
		ShardDoubleHashes:          [][]byte{[]byte(hash)},
		Timestamp:                  expiredTime,
		Status:                     types.Status_STATUS_CHALLENGE_PERIOD,
		Publisher:                  "publisher",
		PublishDataCollateral:      sdk.Coins{},
		SubmitInvalidityCollateral: sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
		PublishedTimestamp:         expiredTime,
	})
	require.NoError(t, err)

	// set invalid status data
	invalidStatusUri := "ipfs://invalid_status"
	err = k.SetPublishedData(ctx, types.PublishedData{
		MetadataUri:                invalidStatusUri,
		ParityShardCount:           0,
		ShardDoubleHashes:          [][]byte{[]byte(hash)},
		Timestamp:                  currentTime,
		Status:                     types.Status_STATUS_CHALLENGING,
		Publisher:                  "publisher",
		PublishDataCollateral:      sdk.Coins{},
		SubmitInvalidityCollateral: sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
		PublishedTimestamp:         currentTime,
	})
	require.NoError(t, err)

	testCases := []struct {
		name      string
		input     *types.MsgSubmitInvalidity
		mockSetup func()
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid sender address",
			input: &types.MsgSubmitInvalidity{
				Sender:      "invalid_address",
				MetadataUri: validMetadataUri,
				Indices:     []int64{0},
			},
			expErr:    true,
			expErrMsg: "invalid sender address",
		},
		{
			name: "indices is not exist",
			input: &types.MsgSubmitInvalidity{
				Sender:      sender.String(),
				MetadataUri: validMetadataUri,
				Indices:     []int64{},
			},
			expErr:    true,
			expErrMsg: "invalid indices",
		},
		{
			name: "invalid metadata uri",
			input: &types.MsgSubmitInvalidity{
				Sender:      sender.String(),
				MetadataUri: "ipfs://nonexistent",
				Indices:     []int64{0},
			},
			expErr:    true,
			expErrMsg: "data not found",
		},
		{
			name: "expired data",
			input: &types.MsgSubmitInvalidity{
				Sender:      sender.String(),
				MetadataUri: expiredMetadataUri,
				Indices:     []int64{0},
			},
			expErr:    true,
			expErrMsg: "challenge period is over",
		},
		{
			name: "invalid status data",
			input: &types.MsgSubmitInvalidity{
				Sender:      sender.String(),
				MetadataUri: invalidStatusUri,
				Indices:     []int64{0},
			},
			expErr:    true,
			expErrMsg: "data is not in challenge period",
		},
		{
			name: "failed to send coins",
			input: &types.MsgSubmitInvalidity{
				Sender:      sender.String(),
				MetadataUri: validMetadataUri,
				Indices:     []int64{0},
			},
			mockSetup: func() {
				// set mock to fail send coins
				mocks.BankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdkerrors.ErrInsufficientFunds)
			},
			expErr:    true,
			expErrMsg: "insufficient funds",
		},
		{
			name: "normal case",
			input: &types.MsgSubmitInvalidity{
				Sender:      sender.String(),
				MetadataUri: validMetadataUri,
				Indices:     []int64{0, 1},
			},
			mockSetup: func() {
				// set mock to success send coins
				mocks.BankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockSetup != nil {
				tc.mockSetup()
			}

			_, err := ms.SubmitInvalidity(ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)

				// in normal case, invalidity is saved in storage
				invalidity, found, err := k.GetInvalidity(ctx, tc.input.MetadataUri, sender)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, tc.input.Sender, invalidity.Sender)
				require.Equal(t, tc.input.MetadataUri, invalidity.MetadataUri)
				require.Equal(t, tc.input.Indices, invalidity.Indices)
			}
		})
	}
}
