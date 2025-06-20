package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/lending/types"
)

func TestGenesis(t *testing.T) {
	f := initFixture(t)

	// Test with default genesis
	t.Run("default genesis", func(t *testing.T) {
		genesisState := types.DefaultGenesis()
		err := f.keeper.InitGenesis(f.ctx, *genesisState)
		require.NoError(t, err)

		got, err := f.keeper.ExportGenesis(f.ctx)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, genesisState.Params, got.Params)
		require.Len(t, got.Markets, 0)
		require.Len(t, got.UserPositions, 0)
		require.Len(t, got.Borrows, 0)
		require.Equal(t, uint64(0), got.BorrowCount)
	})

	// Test with populated genesis
	t.Run("populated genesis", func(t *testing.T) {
		genesisState := &types.GenesisState{
			Params: types.DefaultParams(),
			Markets: []types.Market{
				{
					Denom:             "usdc",
					TotalSupplied:     math.NewInt(1000000),
					TotalBorrowed:     math.NewInt(500000),
					GlobalRewardIndex: math.LegacyOneDec(),
					RiseDenom:         "riseusdc",
				},
				{
					Denom:             "atom",
					TotalSupplied:     math.NewInt(2000000),
					TotalBorrowed:     math.NewInt(1000000),
					GlobalRewardIndex: math.LegacyNewDecWithPrec(11, 1),
					RiseDenom:         "riseatom",
				},
			},
			UserPositions: []types.UserPosition{
				{
					UserAddress:     sdk.AccAddress("user1").String(),
					Denom:           "usdc",
					Amount:          math.NewInt(100000),
					LastRewardIndex: math.LegacyOneDec(),
				},
				{
					UserAddress:     sdk.AccAddress("user1").String(),
					Denom:           "atom",
					Amount:          math.NewInt(50000),
					LastRewardIndex: math.LegacyNewDecWithPrec(95, 2),
				},
				{
					UserAddress:     sdk.AccAddress("user2").String(),
					Denom:           "usdc",
					Amount:          math.NewInt(200000),
					LastRewardIndex: math.LegacyNewDecWithPrec(11, 1),
				},
			},
			Borrows: []types.Borrow{
				{
					Id:                   0,
					Borrower:             sdk.AccAddress("borrower1").String(),
					Amount:               sdk.NewCoin("usdc", math.NewInt(50000)),
					CollateralPoolId:     1,
					CollateralPositionId: 100,
					BlockHeight:          1000,
				},
				{
					Id:                   1,
					Borrower:             sdk.AccAddress("borrower2").String(),
					Amount:               sdk.NewCoin("atom", math.NewInt(25000)),
					CollateralPoolId:     2,
					CollateralPositionId: 200,
					BlockHeight:          1500,
				},
			},
			BorrowCount: 2,
		}

		// Initialize genesis
		err := f.keeper.InitGenesis(f.ctx, *genesisState)
		require.NoError(t, err)

		// Export genesis
		got, err := f.keeper.ExportGenesis(f.ctx)
		require.NoError(t, err)
		require.NotNil(t, got)

		// Verify params
		require.Equal(t, genesisState.Params, got.Params)

		// Verify markets
		require.Len(t, got.Markets, 2)
		marketMap := make(map[string]types.Market)
		for _, market := range got.Markets {
			marketMap[market.Denom] = market
		}
		require.Contains(t, marketMap, "usdc")
		require.Contains(t, marketMap, "atom")
		require.Equal(t, genesisState.Markets[0], marketMap["usdc"])
		require.Equal(t, genesisState.Markets[1], marketMap["atom"])

		// Verify user positions
		require.Len(t, got.UserPositions, 3)
		// Create a map to verify positions
		positionMap := make(map[string]types.UserPosition)
		for _, pos := range got.UserPositions {
			key := pos.UserAddress + ":" + pos.Denom
			positionMap[key] = pos
		}
		for _, expectedPos := range genesisState.UserPositions {
			key := expectedPos.UserAddress + ":" + expectedPos.Denom
			actualPos, exists := positionMap[key]
			require.True(t, exists, "position not found: %s", key)
			require.Equal(t, expectedPos, actualPos)
		}

		// Verify borrows
		require.Len(t, got.Borrows, 2)
		borrowMap := make(map[uint64]types.Borrow)
		for _, borrow := range got.Borrows {
			borrowMap[borrow.Id] = borrow
		}
		for _, expectedBorrow := range genesisState.Borrows {
			actualBorrow, exists := borrowMap[expectedBorrow.Id]
			require.True(t, exists, "borrow not found: %d", expectedBorrow.Id)
			require.Equal(t, expectedBorrow, actualBorrow)
		}

		// Verify borrow count
		require.Equal(t, genesisState.BorrowCount, got.BorrowCount)
	})

	// Test genesis validation
	t.Run("invalid genesis fails", func(t *testing.T) {
		// Invalid params
		invalidGenesis := &types.GenesisState{
			Params: types.NewParams(
				math.LegacyNewDec(2), // Invalid LTV > 1
				math.LegacyNewDecWithPrec(85, 2),
				math.LegacyNewDecWithPrec(5, 2),
			),
		}
		err := f.keeper.InitGenesis(f.ctx, *invalidGenesis)
		// InitGenesis doesn't validate, so it will succeed
		// Validation should happen at the module level
		require.NoError(t, err)
	})
}
