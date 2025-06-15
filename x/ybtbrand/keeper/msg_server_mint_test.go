package keeper_test

import (
	"errors"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/ybtbrand/keeper"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func TestMsgServerMint(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, k keeper.Keeper)
		msg       *types.MsgMint
		setupMock func(mocks moduleMocks)
		wantErr   bool
		errMsg    string
		validate  func(t *testing.T, ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "successful mint",
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
			msg: &types.MsgMint{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Ybt:          sdk.NewCoin(fmt.Sprintf("ybtbase/%s", testAddress2), math.NewInt(1000)),
			},
			setupMock: func(mocks moduleMocks) {
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)
				collateralAddr := keeper.GetCollateralPoolAddress(testAddress)
				brandDenom := keeper.GetTokenDenom(testAddress)
				baseDenom := fmt.Sprintf("ybtbase/%s", testAddress2)

				// Mock transfer base YBT from admin to collateral pool
				baseCoins := sdk.NewCoins(sdk.NewCoin(baseDenom, math.NewInt(1000)))
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), adminAddr, collateralAddr, baseCoins).
					Return(nil)

				// Mock mint brand tokens to admin
				brandCoins := sdk.NewCoins(sdk.NewCoin(brandDenom, math.NewInt(1000)))
				mocks.BankKeeper.EXPECT().
					MintCoins(gomock.Any(), types.ModuleName, brandCoins).
					Return(nil)
				mocks.BankKeeper.EXPECT().
					SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, adminAddr, brandCoins).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "token not found",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Don't create token
			},
			msg: &types.MsgMint{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Ybt:          sdk.NewCoin(fmt.Sprintf("ybtbase/%s", testAddress2), math.NewInt(1000)),
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
			msg: &types.MsgMint{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Ybt:          sdk.NewCoin(fmt.Sprintf("ybtbase/%s", testAddress3), math.NewInt(1000)),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "unauthorized",
		},
		{
			name: "wrong base YBT denom",
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
			msg: &types.MsgMint{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Ybt:          sdk.NewCoin("wrongdenom", math.NewInt(1000)),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "invalid base YBT denom",
		},
		{
			name: "zero amount",
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
			msg: &types.MsgMint{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Ybt:          sdk.NewCoin(fmt.Sprintf("ybtbase/%s", testAddress2), math.ZeroInt()),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "amount must be positive",
		},
		{
			name: "invalid amount - negative",
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
			msg: &types.MsgMint{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Ybt:          sdk.Coin{Denom: fmt.Sprintf("ybtbase/%s", testAddress2), Amount: math.NewInt(-100)},
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "amount must be positive",
		},
		{
			name: "bank transfer fails",
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
			msg: &types.MsgMint{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Ybt:          sdk.NewCoin(fmt.Sprintf("ybtbase/%s", testAddress2), math.NewInt(1000)),
			},
			setupMock: func(mocks moduleMocks) {
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)
				collateralAddr := keeper.GetCollateralPoolAddress(testAddress)
				baseDenom := fmt.Sprintf("ybtbase/%s", testAddress2)

				// Mock transfer base YBT fails
				baseCoins := sdk.NewCoins(sdk.NewCoin(baseDenom, math.NewInt(1000)))
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), adminAddr, collateralAddr, baseCoins).
					Return(errors.New("insufficient funds"))
			},
			wantErr: true,
			errMsg:  "insufficient funds",
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

			resp, err := ms.Mint(ctx, tt.msg)
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
