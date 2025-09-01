package keeper_test

import (
	"errors"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func TestSaveVoteWeightsForBribes(t *testing.T) {
	_, _, addr1 := testdata.KeyTestPubAddr()
	addr1Str := addr1.String()
	_, _, addr2 := testdata.KeyTestPubAddr()
	addr2Str := addr2.String()

	tests := []struct {
		name                string
		votes               []types.Vote
		epochId             uint64
		expectedAllocations map[string]types.BribeAllocation // key: address-epochId-poolId
		expectErr           bool
	}{
		{
			name: "simple case with two voters, two pools",
			votes: []types.Vote{
				{Sender: addr1Str, PoolWeights: []types.PoolWeight{{PoolId: 1, Weight: "1.0"}, {PoolId: 2, Weight: "2.0"}}},
				{Sender: addr2Str, PoolWeights: []types.PoolWeight{{PoolId: 1, Weight: "3.0"}}},
			},
			epochId: 5,
			// Pool 1 total weight: 1.0 + 3.0 = 4.0. Addr1 relative: 1/4=0.25. Addr2 relative: 3/4=0.75
			// Pool 2 total weight: 2.0. Addr1 relative: 2/2=1.0
			expectedAllocations: map[string]types.BribeAllocation{
				addr1Str + "-5-1": {Address: addr1Str, EpochId: 5, PoolId: 1, Weight: "0.250000000000000000", ClaimedBribeIds: []uint64{}},
				addr1Str + "-5-2": {Address: addr1Str, EpochId: 5, PoolId: 2, Weight: "1.000000000000000000", ClaimedBribeIds: []uint64{}},
				addr2Str + "-5-1": {Address: addr2Str, EpochId: 5, PoolId: 1, Weight: "0.750000000000000000", ClaimedBribeIds: []uint64{}},
			},
			expectErr: false,
		},
		{
			name:                "no votes",
			votes:               []types.Vote{},
			epochId:             6,
			expectedAllocations: map[string]types.BribeAllocation{},
			expectErr:           false,
		},
		{
			name: "vote with zero weight",
			votes: []types.Vote{
				{Sender: addr1Str, PoolWeights: []types.PoolWeight{{PoolId: 1, Weight: "1.0"}}},
				{Sender: addr2Str, PoolWeights: []types.PoolWeight{{PoolId: 1, Weight: "0.0"}}},
			},
			epochId: 7,
			expectedAllocations: map[string]types.BribeAllocation{
				addr1Str + "-7-1": {Address: addr1Str, EpochId: 7, PoolId: 1, Weight: "1.000000000000000000", ClaimedBribeIds: []uint64{}},
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fx := initFixture(t)
			sdkCtx := fx.ctx.(sdk.Context)

			// Setup votes
			for _, vote := range tc.votes {
				err := fx.keeper.SetVote(sdkCtx, vote)
				require.NoError(t, err)
			}

			err := fx.keeper.SaveVoteWeightsForBribes(sdkCtx, tc.epochId)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Verify allocations
				allocations, err := fx.keeper.GetAllBribeAllocations(sdkCtx)
				require.NoError(t, err)

				// Create a map from the result for easier lookup
				resultMap := make(map[string]types.BribeAllocation)
				for _, alloc := range allocations {
					key := alloc.Address + "-" + strconv.FormatUint(alloc.EpochId, 10) + "-" + strconv.FormatUint(alloc.PoolId, 10)
					resultMap[key] = alloc
				}

				require.Len(t, resultMap, len(tc.expectedAllocations))
				for key, expected := range tc.expectedAllocations {
					actual, ok := resultMap[key]
					require.True(t, ok, "expected allocation not found: %s", key)
					require.Equal(t, expected.Address, actual.Address)
					require.Equal(t, expected.EpochId, actual.EpochId)
					require.Equal(t, expected.PoolId, actual.PoolId)
					require.Equal(t, expected.Weight, actual.Weight)
				}
			}
		})
	}
}

func TestProcessUnclaimedBribes(t *testing.T) {
	_, _, addr1 := testdata.KeyTestPubAddr()
	addr1Str := addr1.String()
	_, _, addr2 := testdata.KeyTestPubAddr()
	addr2Str := addr2.String()
	feeCollectorAddr := sdk.AccAddress("fee_collector")

	tests := []struct {
		name               string
		epochToProcess     uint64
		initialBribes      []types.Bribe
		initialAllocations []types.BribeAllocation
		setupMocks         func(fx *fixture, expectedUnclaimed sdk.Coins)
		expectErr          bool
	}{
		{
			name:           "fully unclaimed bribe",
			epochToProcess: 4,
			initialBribes: []types.Bribe{
				{Id: 1, EpochId: 4, PoolId: 1, Address: addr1Str, Amount: sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))), ClaimedAmount: sdk.NewCoins()},
			},
			initialAllocations: []types.BribeAllocation{
				{Address: addr2Str, EpochId: 4, PoolId: 1, Weight: "1.0"},
			},
			setupMocks: func(fx *fixture, expectedUnclaimed sdk.Coins) {
				fx.mocks.AcctKeeper.EXPECT().GetModuleAddress("fee_collector").Return(feeCollectorAddr)
				fx.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, feeCollectorAddr, expectedUnclaimed).Return(nil)
			},
			expectErr: false,
		},
		{
			name:           "partially claimed bribe",
			epochToProcess: 4,
			initialBribes: []types.Bribe{
				{Id: 1, EpochId: 4, PoolId: 1, Address: addr1Str, Amount: sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))), ClaimedAmount: sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(40)))},
			},
			initialAllocations: []types.BribeAllocation{
				{Address: addr2Str, EpochId: 4, PoolId: 1, Weight: "1.0"},
			},
			setupMocks: func(fx *fixture, expectedUnclaimed sdk.Coins) {
				fx.mocks.AcctKeeper.EXPECT().GetModuleAddress("fee_collector").Return(feeCollectorAddr)
				fx.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, feeCollectorAddr, expectedUnclaimed).Return(nil)
			},
			expectErr: false,
		},
		{
			name:           "no bribes for epoch",
			epochToProcess: 5,
			setupMocks: func(fx *fixture, expectedUnclaimed sdk.Coins) {
				// No mocks needed as nothing should happen
			},
			expectErr: false,
		},
		{
			name:           "send to fee collector fails",
			epochToProcess: 4,
			initialBribes: []types.Bribe{
				{Id: 1, EpochId: 4, PoolId: 1, Address: addr1Str, Amount: sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))},
			},
			initialAllocations: []types.BribeAllocation{
				{Address: addr2Str, EpochId: 4, PoolId: 1, Weight: "1.0"},
			},
			setupMocks: func(fx *fixture, expectedUnclaimed sdk.Coins) {
				fx.mocks.AcctKeeper.EXPECT().GetModuleAddress("fee_collector").Return(feeCollectorAddr)
				fx.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, feeCollectorAddr, expectedUnclaimed).Return(errors.New("bank send error"))
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fx := initFixture(t)
			sdkCtx := fx.ctx.(sdk.Context)

			// Setup initial state
			for _, bribe := range tc.initialBribes {
				err := fx.keeper.SetBribe(sdkCtx, bribe)
				require.NoError(t, err)
			}
			for _, alloc := range tc.initialAllocations {
				err := fx.keeper.SetBribeAllocation(sdkCtx, alloc)
				require.NoError(t, err)
			}

			// Calculate expected unclaimed amount for the mock
			expectedUnclaimed := sdk.NewCoins()
			for _, bribe := range tc.initialBribes {
				if bribe.EpochId == tc.epochToProcess {
					expectedUnclaimed = expectedUnclaimed.Add(bribe.Amount.Sub(bribe.ClaimedAmount...)...)
				}
			}

			// Setup mocks
			if tc.setupMocks != nil {
				tc.setupMocks(fx, expectedUnclaimed)
			}

			// Process unclaimed bribes
			err := fx.keeper.ProcessUnclaimedBribes(sdkCtx, tc.epochToProcess)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Verify bribes are removed
				bribes, err := fx.keeper.GetAllBribeByEpochId(sdkCtx, tc.epochToProcess)
				require.NoError(t, err)
				require.Len(t, bribes, 0)

				// Verify allocations are removed
				iter, err := fx.keeper.BribeAllocations.Indexes.EpochId.MatchExact(sdkCtx, tc.epochToProcess)
				require.NoError(t, err)
				defer iter.Close()
				require.False(t, iter.Valid(), "allocations for processed epoch should be removed")

				// Verify expired epoch ID is updated
				if len(tc.initialBribes) > 0 || len(tc.initialAllocations) > 0 {
					expiredEpochId := fx.keeper.GetBribeExpiredEpochId(sdkCtx)
					require.Equal(t, tc.epochToProcess, expiredEpochId)
				}
			}
		})
	}
}

func TestFinalizeBribeForEpoch(t *testing.T) {
	feeCollectorAddr := sdk.AccAddress("fee_collector")

	tests := []struct {
		name                 string
		currentEpochId       uint64
		bribeClaimEpochs     uint64
		lastExpiredEpochId   uint64
		setup                func(fx *fixture)
		expectedExpiredEpoch uint64
		expectErr            bool
	}{
		{
			name:               "process two expired epochs",
			currentEpochId:     10,
			bribeClaimEpochs:   5, // Epochs 10-5=5. So, epochs up to 4 should be processed.
			lastExpiredEpochId: 2, // Last processed was 2. So, 3 and 4 should be processed now.
			setup: func(fx *fixture) {
				// Bribes and allocations for epoch 3
				bribe3 := types.Bribe{Id: 1, EpochId: 3, Amount: sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))}
				_ = fx.keeper.SetBribe(fx.ctx, bribe3)
				alloc3 := types.BribeAllocation{Address: "addr1", EpochId: 3, PoolId: 1, Weight: "1.0"}
				_ = fx.keeper.SetBribeAllocation(fx.ctx, alloc3)

				// Bribes and allocations for epoch 4
				bribe4 := types.Bribe{Id: 2, EpochId: 4, Amount: sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(200)))}
				_ = fx.keeper.SetBribe(fx.ctx, bribe4)
				alloc4 := types.BribeAllocation{Address: "addr2", EpochId: 4, PoolId: 2, Weight: "1.0"}
				_ = fx.keeper.SetBribeAllocation(fx.ctx, alloc4)

				// Mocks for processing epoch 3
				fx.mocks.AcctKeeper.EXPECT().GetModuleAddress("fee_collector").Return(feeCollectorAddr)
				fx.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, feeCollectorAddr, bribe3.Amount).Return(nil)

				// Mocks for processing epoch 4
				fx.mocks.AcctKeeper.EXPECT().GetModuleAddress("fee_collector").Return(feeCollectorAddr)
				fx.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, feeCollectorAddr, bribe4.Amount).Return(nil)
			},
			expectedExpiredEpoch: 4,
			expectErr:            false,
		},
		{
			name:                 "no new epochs to process",
			currentEpochId:       10,
			bribeClaimEpochs:     5,
			lastExpiredEpochId:   5, // Already up to date
			setup:                func(fx *fixture) {},
			expectedExpiredEpoch: 5,
			expectErr:            false,
		},
		{
			name:                 "claim period not passed",
			currentEpochId:       4,
			bribeClaimEpochs:     5, // 4 < 5, so no epochs are old enough to be expired
			lastExpiredEpochId:   0,
			setup:                func(fx *fixture) {},
			expectedExpiredEpoch: 0,
			expectErr:            false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fx := initFixture(t)
			sdkCtx := fx.ctx.(sdk.Context)

			// Set initial params and state
			params := types.DefaultParams()
			params.BribeClaimEpochs = tc.bribeClaimEpochs
			err := fx.keeper.Params.Set(sdkCtx, params)
			require.NoError(t, err)
			_ = fx.keeper.SetBribeExpiredEpochId(sdkCtx, tc.lastExpiredEpochId)

			if tc.setup != nil {
				tc.setup(fx)
			}

			// Call the function under test
			err = fx.keeper.FinalizeBribeForEpoch(sdkCtx, tc.currentEpochId)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Check that the expired epoch ID was updated correctly
				finalExpiredId := fx.keeper.GetBribeExpiredEpochId(sdkCtx)
				require.Equal(t, tc.expectedExpiredEpoch, finalExpiredId)
			}
		})
	}
}
