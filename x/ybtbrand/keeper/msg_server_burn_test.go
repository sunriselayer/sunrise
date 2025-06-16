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

func TestMsgServerBurn(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, k keeper.Keeper)
		msg       *types.MsgBurn
		setupMock func(mocks moduleMocks)
		wantErr   bool
		errMsg    string
		validate  func(t *testing.T, ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "successful burn",
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
			msg: &types.MsgBurn{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.NewInt(500),
			},
			setupMock: func(mocks moduleMocks) {
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)
				collateralAddr := keeper.GetCollateralPoolAddress(testAddress)
				brandDenom := keeper.GetTokenDenom(testAddress)
				baseDenom := fmt.Sprintf("ybtbase/%s", testAddress2)

				// Mock burn brand tokens from admin
				brandCoins := sdk.NewCoins(sdk.NewCoin(brandDenom, math.NewInt(500)))
				mocks.BankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), adminAddr, types.ModuleName, brandCoins).
					Return(nil)
				mocks.BankKeeper.EXPECT().
					BurnCoins(gomock.Any(), types.ModuleName, brandCoins).
					Return(nil)

				// Mock transfer base YBT from collateral pool to admin
				baseCoins := sdk.NewCoins(sdk.NewCoin(baseDenom, math.NewInt(500)))
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), collateralAddr, adminAddr, baseCoins).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "token not found",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Don't create token
			},
			msg: &types.MsgBurn{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.NewInt(500),
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
			msg: &types.MsgBurn{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.NewInt(500),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "unauthorized",
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
			msg: &types.MsgBurn{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.ZeroInt(),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "amount must be positive",
		},
		{
			name: "negative amount",
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
			msg: &types.MsgBurn{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.NewInt(-100),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "amount must be positive",
		},
		{
			name: "insufficient balance",
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
			msg: &types.MsgBurn{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.NewInt(500),
			},
			setupMock: func(mocks moduleMocks) {
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)
				brandDenom := keeper.GetTokenDenom(testAddress)

				// Mock burn fails due to insufficient balance
				brandCoins := sdk.NewCoins(sdk.NewCoin(brandDenom, math.NewInt(500)))
				mocks.BankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), adminAddr, types.ModuleName, brandCoins).
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

			resp, err := ms.Burn(ctx, tt.msg)
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
