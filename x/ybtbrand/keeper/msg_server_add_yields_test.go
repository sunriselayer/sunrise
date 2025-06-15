package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/ybtbrand/keeper"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func TestMsgServerAddYields(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, k keeper.Keeper)
		msg       *types.MsgAddYields
		setupMock func(mocks moduleMocks)
		wantErr   bool
		errMsg    string
		validate  func(t *testing.T, ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "successful add yields - multiple denoms",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Initialize yield indices
				err = k.SetYieldIndex(ctx, testAddress, "uusdc", math.LegacyOneDec())
				require.NoError(t, err)
				err = k.SetYieldIndex(ctx, testAddress, "uusdt", math.LegacyOneDec())
				require.NoError(t, err)
			},
			msg: &types.MsgAddYields{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount: sdk.NewCoins(
					sdk.NewCoin("uusdc", math.NewInt(1000)),
					sdk.NewCoin("uusdt", math.NewInt(500)),
				),
			},
			setupMock: func(mocks moduleMocks) {
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)
				yieldPoolAddr1 := keeper.GetYieldPoolAddress(testAddress, "uusdc")
				yieldPoolAddr2 := keeper.GetYieldPoolAddress(testAddress, "uusdt")
				brandDenom := keeper.GetTokenDenom(testAddress)

				// Mock supply for yield index calculation
				mocks.BankKeeper.EXPECT().
					GetSupply(gomock.Any(), brandDenom).
					Return(sdk.NewCoin(brandDenom, math.NewInt(10000)))

				// Mock transfers to yield pools
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), adminAddr, yieldPoolAddr1, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1000)))).
					Return(nil)
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), adminAddr, yieldPoolAddr2, sdk.NewCoins(sdk.NewCoin("uusdt", math.NewInt(500)))).
					Return(nil)

				// Mock set module accounts for yield pools
				mocks.AuthKeeper.EXPECT().
					GetAccount(gomock.Any(), yieldPoolAddr1).
					Return(nil)
				mocks.AuthKeeper.EXPECT().
					SetModuleAccount(gomock.Any(), gomock.Any())
				mocks.AuthKeeper.EXPECT().
					GetAccount(gomock.Any(), yieldPoolAddr2).
					Return(nil)
				mocks.AuthKeeper.EXPECT().
					SetModuleAccount(gomock.Any(), gomock.Any())
			},
			wantErr: false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check yield indices were updated
				// Expected: 1 + (1000 / 10000) = 1.1 for USDC
				// Expected: 1 + (500 / 10000) = 1.05 for USDT
				usdcIndex, found := k.GetYieldIndex(ctx, testAddress, "uusdc")
				require.True(t, found)
				require.Equal(t, math.LegacyNewDecWithPrec(11, 1), usdcIndex)

				usdtIndex, found := k.GetYieldIndex(ctx, testAddress, "uusdt")
				require.True(t, found)
				require.Equal(t, math.LegacyNewDecWithPrec(105, 2), usdtIndex)
			},
		},
		{
			name: "successful add yields - single denom",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Initialize yield index
				err = k.SetYieldIndex(ctx, testAddress, "uusdc", math.LegacyNewDec(2))
				require.NoError(t, err)
			},
			msg: &types.MsgAddYields{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount: sdk.NewCoins(
					sdk.NewCoin("uusdc", math.NewInt(2000)),
				),
			},
			setupMock: func(mocks moduleMocks) {
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)
				yieldPoolAddr := keeper.GetYieldPoolAddress(testAddress, "uusdc")
				brandDenom := keeper.GetTokenDenom(testAddress)

				// Mock supply for yield index calculation
				mocks.BankKeeper.EXPECT().
					GetSupply(gomock.Any(), brandDenom).
					Return(sdk.NewCoin(brandDenom, math.NewInt(5000)))

				// Mock transfer to yield pool
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), adminAddr, yieldPoolAddr, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(2000)))).
					Return(nil)

				// Mock set module account
				mocks.AuthKeeper.EXPECT().
					GetAccount(gomock.Any(), yieldPoolAddr).
					Return(nil)
				mocks.AuthKeeper.EXPECT().
					SetModuleAccount(gomock.Any(), gomock.Any())
			},
			wantErr: false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check yield index was updated
				// Expected: 2 + (2000 / 5000) = 2.4
				usdcIndex, found := k.GetYieldIndex(ctx, testAddress, "uusdc")
				require.True(t, found)
				require.Equal(t, math.LegacyNewDecWithPrec(24, 1), usdcIndex)
			},
		},
		{
			name: "token not found",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Don't create token
			},
			msg: &types.MsgAddYields{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1000))),
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
					Creator:        testAddress,
					Admin:          testAddress2,
					BaseYbtCreator: testAddress3,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgAddYields{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1000))),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "unauthorized",
		},
		{
			name: "empty amount",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgAddYields{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       sdk.NewCoins(),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "amount cannot be empty",
		},
		{
			name: "zero supply - skip yield update",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgAddYields{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1000))),
			},
			setupMock: func(mocks moduleMocks) {
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)
				yieldPoolAddr := keeper.GetYieldPoolAddress(testAddress, "uusdc")
				brandDenom := keeper.GetTokenDenom(testAddress)

				// Mock zero supply
				mocks.BankKeeper.EXPECT().
					GetSupply(gomock.Any(), brandDenom).
					Return(sdk.NewCoin(brandDenom, math.ZeroInt()))

				// Mock transfer to yield pool (still happens even with zero supply)
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), adminAddr, yieldPoolAddr, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1000)))).
					Return(nil)

				// Mock set module account
				mocks.AuthKeeper.EXPECT().
					GetAccount(gomock.Any(), yieldPoolAddr).
					Return(nil)
				mocks.AuthKeeper.EXPECT().
					SetModuleAccount(gomock.Any(), gomock.Any())
			},
			wantErr: false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check yield index remains at default (1.0)
				usdcIndex, found := k.GetYieldIndex(ctx, testAddress, "uusdc")
				require.True(t, found)
				require.Equal(t, math.LegacyOneDec(), usdcIndex)
			},
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

			resp, err := ms.AddYields(ctx, tt.msg)
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
