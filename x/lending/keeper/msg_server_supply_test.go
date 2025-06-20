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

func TestMsgServerSupply(t *testing.T) {
	// Create test addresses
	supplier := sdk.AccAddress("supplier")
	supplierAddr := supplier.String()

	tests := []struct {
		name          string
		msg           *types.MsgSupply
		setupMocks    func(bankKeeper *lendingtest.MockBankKeeper)
		setupKeeper   func(f *fixture)
		wantErr       bool
		errMsg        string
		checkFunc     func(f *fixture)
	}{
		{
			name: "successful supply creates new market and mints rise tokens",
			msg: &types.MsgSupply{
				Sender: supplierAddr,
				Amount: sdk.NewCoin("usdc", math.NewInt(1000000)), // 1 USDC
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// Mock expectations
				bankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), supplier, types.ModuleName, sdk.NewCoins(sdk.NewCoin("usdc", math.NewInt(1000000)))).
					Return(nil)
				
				bankKeeper.EXPECT().
					MintCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(sdk.NewCoin("riseusdc", math.NewInt(1000000)))).
					Return(nil)
				
				bankKeeper.EXPECT().
					SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, supplier, sdk.NewCoins(sdk.NewCoin("riseusdc", math.NewInt(1000000)))).
					Return(nil)
			},
			setupKeeper: func(f *fixture) {
				// No additional keeper setup needed
			},
			wantErr: false,
			checkFunc: func(f *fixture) {
				// Check market was created
				market, err := f.keeper.GetMarket(f.ctx, "usdc")
				require.NoError(t, err)
				require.Equal(t, "usdc", market.Denom)
				require.Equal(t, math.NewInt(1000000), market.TotalSupplied)
				require.Equal(t, math.ZeroInt(), market.TotalBorrowed)
				require.Equal(t, math.LegacyOneDec(), market.GlobalRewardIndex)
				require.Equal(t, "riseusdc", market.RiseDenom)

				// Check user position was created
				position, err := f.keeper.GetUserPosition(f.ctx, supplierAddr, "usdc")
				require.NoError(t, err)
				require.Equal(t, supplierAddr, position.UserAddress)
				require.Equal(t, "usdc", position.Denom)
				require.Equal(t, math.NewInt(1000000), position.Amount) // 1:1 initial ratio
				require.Equal(t, math.LegacyOneDec(), position.LastRewardIndex)
			},
		},
		{
			name: "supply to existing market updates totals",
			msg: &types.MsgSupply{
				Sender: supplierAddr,
				Amount: sdk.NewCoin("usdc", math.NewInt(2000000)), // 2 USDC
			},
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				bankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), supplier, types.ModuleName, sdk.NewCoins(sdk.NewCoin("usdc", math.NewInt(2000000)))).
					Return(nil)
				
				bankKeeper.EXPECT().
					MintCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(sdk.NewCoin("riseusdc", math.NewInt(2000000)))).
					Return(nil)
				
				bankKeeper.EXPECT().
					SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, supplier, sdk.NewCoins(sdk.NewCoin("riseusdc", math.NewInt(2000000)))).
					Return(nil)
			},
			setupKeeper: func(f *fixture) {
				// Create existing market
				market := types.Market{
					Denom:             "usdc",
					TotalSupplied:     math.NewInt(5000000),
					TotalBorrowed:     math.ZeroInt(),
					GlobalRewardIndex: math.LegacyOneDec(),
					RiseDenom:         "riseusdc",
				}
				err := f.keeper.SetMarket(f.ctx, market)
				require.NoError(t, err)
			},
			wantErr: false,
			checkFunc: func(f *fixture) {
				// Check market was updated
				market, err := f.keeper.GetMarket(f.ctx, "usdc")
				require.NoError(t, err)
				require.Equal(t, math.NewInt(7000000), market.TotalSupplied) // 5 + 2 USDC
			},
		},
		{
			name: "invalid sender address",
			msg: &types.MsgSupply{
				Sender: "invalid",
				Amount: sdk.NewCoin("usdc", math.NewInt(1000000)),
			},
			wantErr: true,
			errMsg:  "invalid sender address",
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed for validation error
			},
			setupKeeper: func(f *fixture) {},
		},
		{
			name: "zero amount",
			msg: &types.MsgSupply{
				Sender: supplierAddr,
				Amount: sdk.NewCoin("usdc", math.ZeroInt()),
			},
			wantErr: true,
			errMsg:  "supply amount must be positive",
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed for validation error
			},
			setupKeeper: func(f *fixture) {},
		},
		{
			name: "negative amount",
			msg: &types.MsgSupply{
				Sender: supplierAddr,
				Amount: sdk.Coin{Denom: "usdc", Amount: math.NewInt(-1000)},
			},
			wantErr: true,
			errMsg:  "supply amount must be positive",
			setupMocks: func(bankKeeper *lendingtest.MockBankKeeper) {
				// No mocks needed for validation error
			},
			setupKeeper: func(f *fixture) {},
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
			_, err := msgServer.Supply(f.ctx, tt.msg)

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