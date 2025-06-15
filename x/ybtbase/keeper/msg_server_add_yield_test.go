package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/ybtbase/keeper"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func TestMsgServerAddYield(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, k keeper.Keeper)
		msg       *types.MsgAddYield
		setupMock func(mocks moduleMocks)
		wantErr   bool
		errMsg    string
		validate  func(t *testing.T, ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "successful add yield - first yield",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token first
				token := types.Token{
					Creator:      testAddress,
					Admin:        testAddress,
					Permissioned: false,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
				err = k.SetGlobalRewardIndex(ctx, testAddress, math.LegacyOneDec())
				require.NoError(t, err)
			},
			msg: &types.MsgAddYield{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.NewInt(1000),
			},
			setupMock: func(mocks moduleMocks) {
				denom := keeper.GetTokenDenom(testAddress)
				coins := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(1000)))
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)
				yieldPoolAddr := keeper.GetYieldPoolAddress(testAddress)

				// Mock balance check
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), adminAddr, denom).
					Return(sdk.NewCoin(denom, math.NewInt(5000)))

				// Get total supply
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), gomock.Any(), denom).
					Return(sdk.NewCoin(denom, math.NewInt(10000))).
					AnyTimes()

				// Expect send from admin to yield pool
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), adminAddr, yieldPoolAddr, coins).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check that global reward index was updated
				// If total supply is 10000 and we add 1000 yield, index should increase by 0.1
				index := k.GetGlobalRewardIndex(ctx, testAddress)
				expectedIndex := math.LegacyNewDecWithPrec(11, 1) // 1.1
				require.Equal(t, expectedIndex, index)
			},
		},
		{
			name: "successful add yield - with existing supply",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:      testAddress,
					Admin:        testAddress,
					Permissioned: false,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
				// Set initial index to 1.5
				err = k.SetGlobalRewardIndex(ctx, testAddress, math.LegacyNewDecWithPrec(15, 1))
				require.NoError(t, err)
			},
			msg: &types.MsgAddYield{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.NewInt(2000),
			},
			setupMock: func(mocks moduleMocks) {
				denom := keeper.GetTokenDenom(testAddress)
				coins := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(2000)))
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)
				yieldPoolAddr := keeper.GetYieldPoolAddress(testAddress)

				// Mock balance check
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), adminAddr, denom).
					Return(sdk.NewCoin(denom, math.NewInt(5000)))

				// Get total supply
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), gomock.Any(), denom).
					Return(sdk.NewCoin(denom, math.NewInt(10000))).
					AnyTimes()

				// Expect send from admin to yield pool
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), adminAddr, yieldPoolAddr, coins).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// If total supply is 10000 and we add 2000 yield, index should increase by 0.2
				// From 1.5 to 1.7
				index := k.GetGlobalRewardIndex(ctx, testAddress)
				expectedIndex := math.LegacyNewDecWithPrec(17, 1) // 1.7
				require.Equal(t, expectedIndex, index)
			},
		},
		{
			name: "token not found",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Don't create token
			},
			msg: &types.MsgAddYield{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.NewInt(1000),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "token not found",
		},
		{
			name: "unauthorized - not admin",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token with different admin
				token := types.Token{
					Creator:      testAddress,
					Admin:        testAddress2,
					Permissioned: false,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgAddYield{
				Admin:        testAddress, // Wrong admin
				TokenCreator: testAddress,
				Amount:       math.NewInt(1000),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "unauthorized",
		},
		{
			name: "insufficient balance",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:      testAddress,
					Admin:        testAddress,
					Permissioned: false,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
				err = k.SetGlobalRewardIndex(ctx, testAddress, math.LegacyOneDec())
				require.NoError(t, err)
			},
			msg: &types.MsgAddYield{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.NewInt(6000),
			},
			setupMock: func(mocks moduleMocks) {
				denom := keeper.GetTokenDenom(testAddress)
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)

				// Mock balance check - insufficient balance
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), adminAddr, denom).
					Return(sdk.NewCoin(denom, math.NewInt(5000)))
			},
			wantErr: true,
			errMsg:  "insufficient balance",
		},
		{
			name: "invalid amount - zero",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:      testAddress,
					Admin:        testAddress,
					Permissioned: false,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgAddYield{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.ZeroInt(),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "invalid amount",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := initFixture(t)
			ctx := sdk.UnwrapSDKContext(f.ctx)
			k := f.keeper
			ms := keeper.NewMsgServerImpl(k)

			if tt.setup != nil {
				tt.setup(ctx, k)
			}

			tt.setupMock(f.mocks)

			resp, err := ms.AddYield(ctx, tt.msg)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					require.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)

				if tt.validate != nil {
					tt.validate(t, ctx, k)
				}
			}
		})
	}
}
