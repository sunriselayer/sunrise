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

func TestMsgServerLiquidate(t *testing.T) {
	// Create test addresses
	borrower := sdk.AccAddress("borrower")
	borrowerAddr := borrower.String()
	liquidator := sdk.AccAddress("liquidator")
	liquidatorAddr := liquidator.String()

	tests := []struct {
		name          string
		msg           *types.MsgLiquidate
		setupMocks    func(bankKeeper *lendingtest.MockBankKeeper)
		setupKeeper   func(f *fixture)
		wantErr       bool
		errMsg        string
		checkFunc     func(f *fixture)
	}{
		{
			name: "successful liquidation of undercollateralized position",
			msg: &types.MsgLiquidate{
				Sender:   liquidatorAddr,
				BorrowId: 0,
				Amount:   sdk.NewCoin("usdc", math.NewInt(500000)), // Liquidate half
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// Mock transferring liquidation payment from liquidator to module
				bankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), liquidator, types.ModuleName, sdk.NewCoins(sdk.NewCoin("usdc", math.NewInt(500000)))).
					Return(nil)
				
				// TODO: Mock transferring collateral reward to liquidator
				// This would involve interaction with liquidity pool module
			},
			setupKeeper: func(f *fixture) {
				// Set liquidation threshold to 85%
				params := types.NewParams(
					math.LegacyNewDecWithPrec(80, 2), // 80% LTV
					math.LegacyNewDecWithPrec(85, 2), // 85% liquidation threshold
					math.LegacyNewDecWithPrec(5, 2),  // 5% interest rate
				)
				err := f.keeper.Params.Set(f.ctx, params)
				require.NoError(t, err)
				
				// Create market
				market := types.Market{
					Denom:             "usdc",
					TotalSupplied:     math.NewInt(10000000),
					TotalBorrowed:     math.NewInt(2000000),
					GlobalRewardIndex: math.LegacyOneDec(),
					RiseDenom:         "riseusdc",
				}
				err = f.keeper.SetMarket(f.ctx, market)
				require.NoError(t, err)
				
				// Create undercollateralized borrow
				// Collateral value dropped to $800, but borrowed $1000
				// Health factor = 800 / 1000 = 0.8 < 0.85 (liquidation threshold)
				borrow := types.Borrow{
					Id:                   0,
					Borrower:             borrowerAddr,
					Amount:               sdk.NewCoin("usdc", math.NewInt(1000000)),
					CollateralPoolId:     1,
					CollateralPositionId: 200, // Position with dropped value
					BlockHeight:          100,
				}
				err = f.keeper.Borrows.Set(f.ctx, uint64(0), borrow)
				require.NoError(t, err)
			},
			wantErr: false,
			checkFunc: func(f *fixture) {
				// Check borrow was updated
				borrow, err := f.keeper.Borrows.Get(f.ctx, uint64(0))
				require.NoError(t, err)
				require.Equal(t, math.NewInt(500000), borrow.Amount.Amount) // Half liquidated
				
				// Check market total borrowed was updated
				market, err := f.keeper.GetMarket(f.ctx, "usdc")
				require.NoError(t, err)
				require.Equal(t, math.NewInt(1500000), market.TotalBorrowed) // 2M - 500k
			},
		},
		{
			name: "fail when position is healthy",
			msg: &types.MsgLiquidate{
				Sender:   liquidatorAddr,
				BorrowId: 0,
				Amount:   sdk.NewCoin("usdc", math.NewInt(500000)),
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed - should fail before bank operations
			},
			setupKeeper: func(f *fixture) {
				// Create healthy borrow
				// Collateral value $2000, borrowed $1000
				// Health factor = 2000 / 1000 = 2.0 > 0.85 (liquidation threshold)
				borrow := types.Borrow{
					Id:                   0,
					Borrower:             borrowerAddr,
					Amount:               sdk.NewCoin("usdc", math.NewInt(1000000)),
					CollateralPoolId:     1,
					CollateralPositionId: 100, // Position with good collateral
					BlockHeight:          100,
				}
				err := f.keeper.Borrows.Set(f.ctx, uint64(0), borrow)
				require.NoError(t, err)
				
				// Create market
				market := types.Market{
					Denom:             "usdc",
					TotalSupplied:     math.NewInt(10000000),
					TotalBorrowed:     math.NewInt(1000000),
					GlobalRewardIndex: math.LegacyOneDec(),
					RiseDenom:         "riseusdc",
				}
				err = f.keeper.SetMarket(f.ctx, market)
				require.NoError(t, err)
			},
			wantErr: true,
			errMsg:  "position is not undercollateralized",
		},
		{
			name: "fail when borrow not found",
			msg: &types.MsgLiquidate{
				Sender:   liquidatorAddr,
				BorrowId: 999,
				Amount:   sdk.NewCoin("usdc", math.NewInt(500000)),
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed - should fail before bank operations
			},
			setupKeeper: func(f *fixture) {
				// No setup needed
			},
			wantErr: true,
			errMsg:  "borrow not found",
		},
		{
			name: "fail when wrong denom",
			msg: &types.MsgLiquidate{
				Sender:   liquidatorAddr,
				BorrowId: 0,
				Amount:   sdk.NewCoin("atom", math.NewInt(500000)), // Wrong denom
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed - should fail before bank operations
			},
			setupKeeper: func(f *fixture) {
				// Create borrow
				borrow := types.Borrow{
					Id:                   0,
					Borrower:             borrowerAddr,
					Amount:               sdk.NewCoin("usdc", math.NewInt(1000000)),
					CollateralPoolId:     1,
					CollateralPositionId: 200,
					BlockHeight:          100,
				}
				err := f.keeper.Borrows.Set(f.ctx, uint64(0), borrow)
				require.NoError(t, err)
			},
			wantErr: true,
			errMsg:  "liquidation denom mismatch",
		},
		{
			name: "fail when liquidating too much",
			msg: &types.MsgLiquidate{
				Sender:   liquidatorAddr,
				BorrowId: 0,
				Amount:   sdk.NewCoin("usdc", math.NewInt(2000000)), // More than borrowed
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed - should fail before bank operations
			},
			setupKeeper: func(f *fixture) {
				// Create borrow
				borrow := types.Borrow{
					Id:                   0,
					Borrower:             borrowerAddr,
					Amount:               sdk.NewCoin("usdc", math.NewInt(1000000)),
					CollateralPoolId:     1,
					CollateralPositionId: 200,
					BlockHeight:          100,
				}
				err := f.keeper.Borrows.Set(f.ctx, uint64(0), borrow)
				require.NoError(t, err)
			},
			wantErr: true,
			errMsg:  "liquidation amount exceeds debt",
		},
		{
			name: "invalid sender address",
			msg: &types.MsgLiquidate{
				Sender:   "invalid",
				BorrowId: 0,
				Amount:   sdk.NewCoin("usdc", math.NewInt(500000)),
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed for validation error
			},
			setupKeeper: func(f *fixture) {},
			wantErr:     true,
			errMsg:      "invalid sender address",
		},
		{
			name: "zero amount",
			msg: &types.MsgLiquidate{
				Sender:   liquidatorAddr,
				BorrowId: 0,
				Amount:   sdk.NewCoin("usdc", math.ZeroInt()),
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed for validation error
			},
			setupKeeper: func(f *fixture) {},
			wantErr:     true,
			errMsg:      "liquidation amount must be positive",
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
			_, err := msgServer.Liquidate(f.ctx, tt.msg)

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