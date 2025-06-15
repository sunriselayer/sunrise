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

func TestMsgServerClaimYields(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, k keeper.Keeper)
		msg       *types.MsgClaimYields
		setupMock func(mocks moduleMocks)
		wantErr   bool
		errMsg    string
		validate  func(t *testing.T, ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "successful claim yields - multiple denoms",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Set yield indices
				err = k.SetYieldIndex(ctx, testAddress, "uusdc", math.LegacyNewDecWithPrec(15, 1)) // 1.5
				require.NoError(t, err)
				err = k.SetYieldIndex(ctx, testAddress, "uusdt", math.LegacyNewDec(2)) // 2.0
				require.NoError(t, err)

				// Set user's last yield indices
				err = k.SetUserLastYieldIndex(ctx, testAddress, testAddress2, "uusdc", math.LegacyOneDec()) // 1.0
				require.NoError(t, err)
				err = k.SetUserLastYieldIndex(ctx, testAddress, testAddress2, "uusdt", math.LegacyOneDec()) // 1.0
				require.NoError(t, err)
			},
			msg: &types.MsgClaimYields{
				Sender:       testAddress2,
				TokenCreator: testAddress,
				Denoms:       []string{"uusdc", "uusdt"},
			},
			setupMock: func(mocks moduleMocks) {
				senderAddr := sdk.MustAccAddressFromBech32(testAddress2)
				yieldPoolAddr1 := keeper.GetYieldPoolAddress(testAddress, "uusdc")
				yieldPoolAddr2 := keeper.GetYieldPoolAddress(testAddress, "uusdt")
				brandDenom := keeper.GetTokenDenom(testAddress)

				// Mock user balance - 1000 brand tokens
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), senderAddr, brandDenom).
					Return(sdk.NewCoin(brandDenom, math.NewInt(1000)))

				// Expected yield USDC: 1000 * (1.5 - 1.0) = 500
				// Expected yield USDT: 1000 * (2.0 - 1.0) = 1000

				// Mock yield pool balances
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), yieldPoolAddr1, "uusdc").
					Return(sdk.NewCoin("uusdc", math.NewInt(2000)))
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), yieldPoolAddr2, "uusdt").
					Return(sdk.NewCoin("uusdt", math.NewInt(3000)))

				// Mock transfers from yield pools to user
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), yieldPoolAddr1, senderAddr, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(500)))).
					Return(nil)
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), yieldPoolAddr2, senderAddr, sdk.NewCoins(sdk.NewCoin("uusdt", math.NewInt(1000)))).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check user's last yield indices were updated
				usdcIndex, found := k.GetUserLastYieldIndex(ctx, testAddress, testAddress2, "uusdc")
				require.True(t, found)
				require.Equal(t, math.LegacyNewDecWithPrec(15, 1), usdcIndex)
				
				usdtIndex, found := k.GetUserLastYieldIndex(ctx, testAddress, testAddress2, "uusdt")
				require.True(t, found)
				require.Equal(t, math.LegacyNewDec(2), usdtIndex)
			},
		},
		{
			name: "successful claim yields - single denom",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Set yield index
				err = k.SetYieldIndex(ctx, testAddress, "uusdc", math.LegacyNewDec(3))
				require.NoError(t, err)

				// Set user's last yield index
				err = k.SetUserLastYieldIndex(ctx, testAddress, testAddress2, "uusdc", math.LegacyNewDec(2))
				require.NoError(t, err)
			},
			msg: &types.MsgClaimYields{
				Sender:       testAddress2,
				TokenCreator: testAddress,
				Denoms:       []string{"uusdc"},
			},
			setupMock: func(mocks moduleMocks) {
				senderAddr := sdk.MustAccAddressFromBech32(testAddress2)
				yieldPoolAddr := keeper.GetYieldPoolAddress(testAddress, "uusdc")
				brandDenom := keeper.GetTokenDenom(testAddress)

				// Mock user balance
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), senderAddr, brandDenom).
					Return(sdk.NewCoin(brandDenom, math.NewInt(500)))

				// Expected yield: 500 * (3.0 - 2.0) = 500

				// Mock yield pool balance
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), yieldPoolAddr, "uusdc").
					Return(sdk.NewCoin("uusdc", math.NewInt(1000)))

				// Mock transfer
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), yieldPoolAddr, senderAddr, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(500)))).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check user's last yield index was updated
				usdcIndex, found := k.GetUserLastYieldIndex(ctx, testAddress, testAddress2, "uusdc")
				require.True(t, found)
				require.Equal(t, math.LegacyNewDec(3), usdcIndex)
			},
		},
		{
			name: "no yield to claim - indices already up to date",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Set both indices to same value
				err = k.SetYieldIndex(ctx, testAddress, "uusdc", math.LegacyNewDec(2))
				require.NoError(t, err)
				err = k.SetUserLastYieldIndex(ctx, testAddress, testAddress2, "uusdc", math.LegacyNewDec(2))
				require.NoError(t, err)
			},
			msg: &types.MsgClaimYields{
				Sender:       testAddress2,
				TokenCreator: testAddress,
				Denoms:       []string{"uusdc"},
			},
			setupMock: func(mocks moduleMocks) {
				senderAddr := sdk.MustAccAddressFromBech32(testAddress2)
				brandDenom := keeper.GetTokenDenom(testAddress)

				// Mock user balance
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), senderAddr, brandDenom).
					Return(sdk.NewCoin(brandDenom, math.NewInt(1000)))
			},
			wantErr: true,
			errMsg:  "no yield to claim",
		},
		{
			name: "no balance - nothing to claim",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Set yield indices
				err = k.SetYieldIndex(ctx, testAddress, "uusdc", math.LegacyNewDec(2))
				require.NoError(t, err)
			},
			msg: &types.MsgClaimYields{
				Sender:       testAddress2,
				TokenCreator: testAddress,
				Denoms:       []string{"uusdc"},
			},
			setupMock: func(mocks moduleMocks) {
				senderAddr := sdk.MustAccAddressFromBech32(testAddress2)
				brandDenom := keeper.GetTokenDenom(testAddress)

				// Mock zero balance
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), senderAddr, brandDenom).
					Return(sdk.NewCoin(brandDenom, math.ZeroInt()))
			},
			wantErr: true,
			errMsg:  "no balance",
		},
		{
			name: "insufficient yield pool balance",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Set yield indices
				err = k.SetYieldIndex(ctx, testAddress, "uusdc", math.LegacyNewDec(2))
				require.NoError(t, err)
				err = k.SetUserLastYieldIndex(ctx, testAddress, testAddress2, "uusdc", math.LegacyOneDec())
				require.NoError(t, err)
			},
			msg: &types.MsgClaimYields{
				Sender:       testAddress2,
				TokenCreator: testAddress,
				Denoms:       []string{"uusdc"},
			},
			setupMock: func(mocks moduleMocks) {
				senderAddr := sdk.MustAccAddressFromBech32(testAddress2)
				yieldPoolAddr := keeper.GetYieldPoolAddress(testAddress, "uusdc")
				brandDenom := keeper.GetTokenDenom(testAddress)

				// Mock user balance
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), senderAddr, brandDenom).
					Return(sdk.NewCoin(brandDenom, math.NewInt(1000)))

				// Mock insufficient yield pool balance
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), yieldPoolAddr, "uusdc").
					Return(sdk.NewCoin("uusdc", math.NewInt(100)))
			},
			wantErr: true,
			errMsg:  "insufficient yield pool balance",
		},
		{
			name: "token not found",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Don't create token
			},
			msg: &types.MsgClaimYields{
				Sender:       testAddress2,
				TokenCreator: testAddress,
				Denoms:       []string{"uusdc"},
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "token not found",
		},
		{
			name: "empty denoms",
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
			msg: &types.MsgClaimYields{
				Sender:       testAddress2,
				TokenCreator: testAddress,
				Denoms:       []string{},
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "denoms cannot be empty",
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

			resp, err := ms.ClaimYields(ctx, tt.msg)
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