package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/lending/keeper"
	lendingtest "github.com/sunriselayer/sunrise/x/lending/testutil"
	"github.com/sunriselayer/sunrise/x/lending/types"
)

func TestMsgServerBorrow(t *testing.T) {
	// Create test addresses
	borrower := sdk.AccAddress("borrower")
	borrowerAddr := borrower.String()

	tests := []struct {
		name          string
		msg           *types.MsgBorrow
		setupMocks    func(bankKeeper *lendingtest.MockBankKeeper)
		setupKeeper   func(f *fixture)
		wantErr       bool
		errMsg        string
		checkFunc     func(f *fixture)
	}{
		{
			name: "successful borrow with sufficient collateral",
			msg: &types.MsgBorrow{
				Sender:               borrowerAddr,
				Amount:               sdk.NewCoin("usdc", math.NewInt(1000000)), // 1 USDC
				CollateralPoolId:     1,
				CollateralPositionId: 100,
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// Mock sending borrowed tokens to borrower
				bankKeeper.EXPECT().
					SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, borrower, sdk.NewCoins(sdk.NewCoin("usdc", math.NewInt(1000000)))).
					Return(nil)
			},
			setupKeeper: func(f *fixture) {
				// Create market with available liquidity
				market := types.Market{
					Denom:             "usdc",
					TotalSupplied:     math.NewInt(10000000), // 10 USDC
					TotalBorrowed:     math.NewInt(2000000),  // 2 USDC already borrowed
					GlobalRewardIndex: math.LegacyOneDec(),
					RiseDenom:         "riseusdc",
				}
				err := f.keeper.SetMarket(f.ctx, market)
				require.NoError(t, err)
				
				// TODO: Mock TWAP price oracle
				// For now, assume 1 USDC = 1 USD, collateral value = 2000000 USD
				// With 80% LTV, can borrow up to 1600000 USDC
				// Borrowing 1000000 USDC should be safe
			},
			wantErr: false,
			checkFunc: func(f *fixture) {
				// Check borrow was created
				borrow, err := f.keeper.Borrows.Get(f.ctx, uint64(0))
				require.NoError(t, err)
				require.Equal(t, borrowerAddr, borrow.Borrower)
				require.Equal(t, "usdc", borrow.Amount.Denom)
				require.Equal(t, math.NewInt(1000000), borrow.Amount.Amount)
				require.Equal(t, uint64(1), borrow.CollateralPoolId)
				require.Equal(t, uint64(100), borrow.CollateralPositionId)
				
				// Check market total borrowed was updated
				market, err := f.keeper.GetMarket(f.ctx, "usdc")
				require.NoError(t, err)
				require.Equal(t, math.NewInt(3000000), market.TotalBorrowed) // 2 + 1 USDC
				
				// Check borrow ID was incremented
				nextId, err := f.keeper.BorrowId.Peek(f.ctx)
				require.NoError(t, err)
				require.Equal(t, uint64(1), nextId)
			},
		},
		{
			name: "fail when market does not exist",
			msg: &types.MsgBorrow{
				Sender:               borrowerAddr,
				Amount:               sdk.NewCoin("nonexistent", math.NewInt(1000000)),
				CollateralPoolId:     1,
				CollateralPositionId: 100,
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed - should fail before bank operations
			},
			setupKeeper: func(f *fixture) {
				// No setup needed
			},
			wantErr: true,
			errMsg:  "market not found",
		},
		{
			name: "fail when insufficient liquidity",
			msg: &types.MsgBorrow{
				Sender:               borrowerAddr,
				Amount:               sdk.NewCoin("usdc", math.NewInt(1000000)), // Trying to borrow 1 USDC
				CollateralPoolId:     1,
				CollateralPositionId: 100,
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed - should fail before bank operations
			},
			setupKeeper: func(f *fixture) {
				// Create market with no available liquidity
				market := types.Market{
					Denom:             "usdc",
					TotalSupplied:     math.NewInt(1000000), // 1 USDC
					TotalBorrowed:     math.NewInt(1000000), // 1 USDC (all borrowed)
					GlobalRewardIndex: math.LegacyOneDec(),
					RiseDenom:         "riseusdc",
				}
				err := f.keeper.SetMarket(f.ctx, market)
				require.NoError(t, err)
			},
			wantErr: true,
			errMsg:  "insufficient liquidity",
		},
		{
			name: "fail when insufficient collateral",
			msg: &types.MsgBorrow{
				Sender:               borrowerAddr,
				Amount:               sdk.NewCoin("usdc", math.NewInt(1000000)), // Trying to borrow 1 USDC
				CollateralPoolId:     1,
				CollateralPositionId: 101, // Position with low collateral
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed - should fail before bank operations
			},
			setupKeeper: func(f *fixture) {
				// Create market
				market := types.Market{
					Denom:             "usdc",
					TotalSupplied:     math.NewInt(10000000), // 10 USDC
					TotalBorrowed:     math.NewInt(0),
					GlobalRewardIndex: math.LegacyOneDec(),
					RiseDenom:         "riseusdc",
				}
				err := f.keeper.SetMarket(f.ctx, market)
				require.NoError(t, err)
				
				// TODO: Mock TWAP price oracle to return low collateral value
				// For now, assume collateral value = 100000 USD
				// With 80% LTV, can only borrow up to 80000 USDC
				// Trying to borrow 1000000 USDC should fail
			},
			wantErr: true,
			errMsg:  "insufficient collateral",
		},
		{
			name: "invalid sender address",
			msg: &types.MsgBorrow{
				Sender:               "invalid",
				Amount:               sdk.NewCoin("usdc", math.NewInt(1000000)),
				CollateralPoolId:     1,
				CollateralPositionId: 100,
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed for validation error
			},
			setupKeeper: func(f *fixture) {},
			wantErr:     true,
			errMsg:      "invalid sender address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset keeper state for each test
			f := initFixture(t)
			defer f.ctrl.Finish()
			
			// Setup mocks
			if tt.setupMocks != nil {
				tt.setupMocks(f.bankKeeper)
			}
			
			// Setup keeper state
			if tt.setupKeeper != nil {
				tt.setupKeeper(f)
			}
			
			msgServer := keeper.NewMsgServerImpl(f.keeper)

			// Execute
			_, err := msgServer.Borrow(f.ctx, tt.msg)

			// Check error
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				// Run additional checks
				if tt.checkFunc != nil {
					tt.checkFunc(f)
				}
			}
		})
	}
}