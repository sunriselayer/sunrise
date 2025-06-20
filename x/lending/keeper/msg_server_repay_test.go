package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/lending/keeper"
	lendingtest "github.com/sunriselayer/sunrise/x/lending/testutil"
	"github.com/sunriselayer/sunrise/x/lending/types"
)

func TestMsgServerRepay(t *testing.T) {
	// Create test addresses
	borrower := sdk.AccAddress("borrower")
	borrowerAddr := borrower.String()
	otherUser := sdk.AccAddress("otheruser")
	otherUserAddr := otherUser.String()

	tests := []struct {
		name          string
		msg           *types.MsgRepay
		setupMocks    func(bankKeeper *lendingtest.MockBankKeeper)
		setupKeeper   func(f *fixture)
		wantErr       bool
		errMsg        string
		checkFunc     func(f *fixture)
	}{
		{
			name: "successful full repayment",
			msg: &types.MsgRepay{
				Sender:   borrowerAddr,
				BorrowId: 0,
				Amount:   sdk.NewCoin("usdc", math.NewInt(1000000)), // Full repayment
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// Mock transferring repayment from user to module
				bankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), borrower, types.ModuleName, sdk.NewCoins(sdk.NewCoin("usdc", math.NewInt(1000000)))).
					Return(nil)
			},
			setupKeeper: func(f *fixture) {
				// Create market
				market := types.Market{
					Denom:             "usdc",
					TotalSupplied:     math.NewInt(10000000),
					TotalBorrowed:     math.NewInt(3000000), // 3 USDC borrowed
					GlobalRewardIndex: math.LegacyOneDec(),
					RiseDenom:         "riseusdc",
				}
				err := f.keeper.SetMarket(f.ctx, market)
				require.NoError(t, err)
				
				// Create borrow
				borrow := types.Borrow{
					Id:                   0,
					Borrower:             borrowerAddr,
					Amount:               sdk.NewCoin("usdc", math.NewInt(1000000)),
					CollateralPoolId:     1,
					CollateralPositionId: 100,
					BlockHeight:          100,
				}
				err = f.keeper.Borrows.Set(f.ctx, uint64(0), borrow)
				require.NoError(t, err)
			},
			wantErr: false,
			checkFunc: func(f *fixture) {
				// Check borrow was deleted
				_, err := f.keeper.Borrows.Get(f.ctx, uint64(0))
				require.Error(t, err)
				require.ErrorIs(t, err, collections.ErrNotFound)
				
				// Check market total borrowed was updated
				market, err := f.keeper.GetMarket(f.ctx, "usdc")
				require.NoError(t, err)
				require.Equal(t, math.NewInt(2000000), market.TotalBorrowed) // 3 - 1 USDC
			},
		},
		{
			name: "successful partial repayment",
			msg: &types.MsgRepay{
				Sender:   borrowerAddr,
				BorrowId: 0,
				Amount:   sdk.NewCoin("usdc", math.NewInt(500000)), // Partial repayment
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// Mock transferring repayment from user to module
				bankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), borrower, types.ModuleName, sdk.NewCoins(sdk.NewCoin("usdc", math.NewInt(500000)))).
					Return(nil)
			},
			setupKeeper: func(f *fixture) {
				// Create market
				market := types.Market{
					Denom:             "usdc",
					TotalSupplied:     math.NewInt(10000000),
					TotalBorrowed:     math.NewInt(3000000),
					GlobalRewardIndex: math.LegacyOneDec(),
					RiseDenom:         "riseusdc",
				}
				err := f.keeper.SetMarket(f.ctx, market)
				require.NoError(t, err)
				
				// Create borrow
				borrow := types.Borrow{
					Id:                   0,
					Borrower:             borrowerAddr,
					Amount:               sdk.NewCoin("usdc", math.NewInt(1000000)),
					CollateralPoolId:     1,
					CollateralPositionId: 100,
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
				require.Equal(t, math.NewInt(500000), borrow.Amount.Amount) // 1M - 500k
				
				// Check market total borrowed was updated
				market, err := f.keeper.GetMarket(f.ctx, "usdc")
				require.NoError(t, err)
				require.Equal(t, math.NewInt(2500000), market.TotalBorrowed) // 3M - 500k
			},
		},
		{
			name: "fail when borrow not found",
			msg: &types.MsgRepay{
				Sender:   borrowerAddr,
				BorrowId: 999,
				Amount:   sdk.NewCoin("usdc", math.NewInt(1000000)),
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
			name: "fail when not borrower",
			msg: &types.MsgRepay{
				Sender:   otherUserAddr,
				BorrowId: 0,
				Amount:   sdk.NewCoin("usdc", math.NewInt(1000000)),
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed - should fail before bank operations
			},
			setupKeeper: func(f *fixture) {
				// Create borrow with different borrower
				borrow := types.Borrow{
					Id:                   0,
					Borrower:             borrowerAddr,
					Amount:               sdk.NewCoin("usdc", math.NewInt(1000000)),
					CollateralPoolId:     1,
					CollateralPositionId: 100,
					BlockHeight:          100,
				}
				err := f.keeper.Borrows.Set(f.ctx, uint64(0), borrow)
				require.NoError(t, err)
			},
			wantErr: true,
			errMsg:  "unauthorized",
		},
		{
			name: "fail when wrong denom",
			msg: &types.MsgRepay{
				Sender:   borrowerAddr,
				BorrowId: 0,
				Amount:   sdk.NewCoin("atom", math.NewInt(1000000)), // Wrong denom
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
					CollateralPositionId: 100,
					BlockHeight:          100,
				}
				err := f.keeper.Borrows.Set(f.ctx, uint64(0), borrow)
				require.NoError(t, err)
			},
			wantErr: true,
			errMsg:  "repayment denom mismatch",
		},
		{
			name: "fail when overpaying",
			msg: &types.MsgRepay{
				Sender:   borrowerAddr,
				BorrowId: 0,
				Amount:   sdk.NewCoin("usdc", math.NewInt(2000000)), // Trying to pay more than owed
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
					CollateralPositionId: 100,
					BlockHeight:          100,
				}
				err := f.keeper.Borrows.Set(f.ctx, uint64(0), borrow)
				require.NoError(t, err)
			},
			wantErr: true,
			errMsg:  "repayment exceeds debt",
		},
		{
			name: "invalid sender address",
			msg: &types.MsgRepay{
				Sender:   "invalid",
				BorrowId: 0,
				Amount:   sdk.NewCoin("usdc", math.NewInt(1000000)),
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
			msg: &types.MsgRepay{
				Sender:   borrowerAddr,
				BorrowId: 0,
				Amount:   sdk.NewCoin("usdc", math.ZeroInt()),
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed for validation error
			},
			setupKeeper: func(f *fixture) {},
			wantErr:     true,
			errMsg:      "repayment amount must be positive",
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
			_, err := msgServer.Repay(f.ctx, tt.msg)

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