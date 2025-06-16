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

func TestMsgServerMint(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, k keeper.Keeper)
		msg       *types.MsgMint
		setupMock func(mocks moduleMocks)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful mint",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token first
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
				err = k.SetGlobalRewardIndex(ctx, testAddress, math.LegacyOneDec())
				require.NoError(t, err)
			},
			msg: &types.MsgMint{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.NewInt(1000),
			},
			setupMock: func(mocks moduleMocks) {
				denom := keeper.GetTokenDenom(testAddress)
				coins := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(1000)))

				// Expect mint coins to module account
				mocks.BankKeeper.EXPECT().
					MintCoins(gomock.Any(), types.ModuleName, coins).
					Return(nil)

				// Expect send from module to admin
				mocks.BankKeeper.EXPECT().
					SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, sdk.MustAccAddressFromBech32(testAddress), coins).
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
					Creator:        testAddress,
					Admin:          testAddress2,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgMint{
				Admin:        testAddress, // Wrong admin
				TokenCreator: testAddress,
				Amount:       math.NewInt(1000),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "unauthorized",
		},
		{
			name: "invalid amount - zero",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgMint{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.ZeroInt(),
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "invalid amount",
		},
		{
			name: "invalid amount - negative",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgMint{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Amount:       math.NewInt(-1000),
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

			resp, err := ms.Mint(ctx, tt.msg)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					require.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}
