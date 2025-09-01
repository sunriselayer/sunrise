package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/da/keeper"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestGetZkpThreshold(t *testing.T) {
	// Mock validator
	pubKey := ed25519.GenPrivKey().PubKey()
	anyPk, err := codectypes.NewAnyWithValue(pubKey)
	require.NoError(t, err)
	mockVal := &stakingtypes.Validator{
		OperatorAddress:   sdk.ValAddress("validator").String(),
		ConsensusPubkey:   anyPk,
		Status:            stakingtypes.Bonded,
		Tokens:            math.NewInt(100),
		DelegatorShares:   math.LegacyNewDec(100),
		MinSelfDelegation: math.NewInt(1),
	}

	testCases := []struct {
		name                string
		setupMocks          func(mocks DaMocks, k keeper.Keeper, sdkCtx sdk.Context)
		shardCount          uint64
		expectedThreshold   uint64
		expectedErr         error
		replicationFactor   string
		numActiveValidators int
	}{
		{
			name:                "normal case",
			shardCount:          100,
			expectedThreshold:   20,
			replicationFactor:   "2.0",
			numActiveValidators: 10,
			setupMocks: func(mocks DaMocks, k keeper.Keeper, sdkCtx sdk.Context) {
				params := types.DefaultParams()
				params.ReplicationFactor = "2.0"
				require.NoError(t, k.Params.Set(sdkCtx, params))

				valAddrs := make([]sdk.ValAddress, 10)
				for i := 0; i < 10; i++ {
					valAddrs[i] = sdk.ValAddress(ed25519.GenPrivKey().PubKey().Address())
				}

				iterator := &mockIterator{valAddrs: valAddrs}
				mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(iterator, nil)
				for _, valAddr := range valAddrs {
					mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr).Return(mockVal, nil)
				}
			},
		},
		{
			name:                "zero active validators",
			shardCount:          100,
			expectedThreshold:   0,
			replicationFactor:   "2.0",
			numActiveValidators: 0,
			setupMocks: func(mocks DaMocks, k keeper.Keeper, sdkCtx sdk.Context) {
				params := types.DefaultParams()
				params.ReplicationFactor = "2.0"
				require.NoError(t, k.Params.Set(sdkCtx, params))

				iterator := &mockIterator{valAddrs: []sdk.ValAddress{}}
				mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(iterator, nil)
			},
		},
		{
			name:                "threshold exceeds shard count",
			shardCount:          10,
			expectedThreshold:   10,
			replicationFactor:   "3.0",
			numActiveValidators: 2,
			setupMocks: func(mocks DaMocks, k keeper.Keeper, sdkCtx sdk.Context) {
				params := types.DefaultParams()
				params.ReplicationFactor = "3.0"
				require.NoError(t, k.Params.Set(sdkCtx, params))

				valAddrs := make([]sdk.ValAddress, 2)
				for i := 0; i < 2; i++ {
					valAddrs[i] = sdk.ValAddress(ed25519.GenPrivKey().PubKey().Address())
				}

				iterator := &mockIterator{valAddrs: valAddrs}
				mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(iterator, nil)
				for _, valAddr := range valAddrs {
					mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr).Return(mockVal, nil)
				}
			},
		},
		{
			name:                "threshold less than 1",
			shardCount:          10,
			expectedThreshold:   1,
			replicationFactor:   "0.1",
			numActiveValidators: 10,
			setupMocks: func(mocks DaMocks, k keeper.Keeper, sdkCtx sdk.Context) {
				params := types.DefaultParams()
				params.ReplicationFactor = "0.1"
				require.NoError(t, k.Params.Set(sdkCtx, params))

				valAddrs := make([]sdk.ValAddress, 10)
				for i := 0; i < 10; i++ {
					valAddrs[i] = sdk.ValAddress(ed25519.GenPrivKey().PubKey().Address())
				}

				iterator := &mockIterator{valAddrs: valAddrs}
				mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(iterator, nil)
				for _, valAddr := range valAddrs {
					mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr).Return(mockVal, nil)
				}
			},
		},
		{
			name:        "iterator error",
			shardCount:  100,
			expectedErr: errors.New("iterator error"),
			setupMocks: func(mocks DaMocks, k keeper.Keeper, sdkCtx sdk.Context) {
				mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(nil, errors.New("iterator error"))
			},
		},
		{
			name:              "validator error",
			shardCount:        100,
			expectedErr:       nil, // error is logged and skipped
			expectedThreshold: 0,
			setupMocks: func(mocks DaMocks, k keeper.Keeper, sdkCtx sdk.Context) {
				params := types.DefaultParams()
				params.ReplicationFactor = "2.0"
				require.NoError(t, k.Params.Set(sdkCtx, params))

				valAddrs := []sdk.ValAddress{sdk.ValAddress("validator")}
				iterator := &mockIterator{valAddrs: valAddrs}
				mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(iterator, nil)
				mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddrs[0]).Return(nil, errors.New("validator error"))
			},
		},
		{
			name:        "params error - invalid replication factor",
			shardCount:  100,
			expectedErr: errors.New("decimal string cannot be empty"),
			setupMocks: func(mocks DaMocks, k keeper.Keeper, sdkCtx sdk.Context) {
				params := types.DefaultParams()
				params.ReplicationFactor = "" // Set invalid param
				require.NoError(t, k.Params.Set(sdkCtx, params))

				iterator := &mockIterator{valAddrs: []sdk.ValAddress{sdk.ValAddress("validator")}}
				mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(iterator, nil)
				mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), gomock.Any()).Return(mockVal, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k, mocks, _, ctx := setupMsgServer(t)
			sdkCtx := sdk.UnwrapSDKContext(ctx)

			if tc.setupMocks != nil {
				tc.setupMocks(mocks, k, sdkCtx)
			}

			threshold, err := k.GetZkpThreshold(sdkCtx, tc.shardCount)

			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedThreshold, threshold)
			}
		})
	}
}
