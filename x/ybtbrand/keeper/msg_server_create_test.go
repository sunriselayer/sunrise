package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/ybtbrand/keeper"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
	ybtbasetypes "github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func TestMsgServerCreate(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, k keeper.Keeper)
		msg       *types.MsgCreate
		setupMock func(mocks moduleMocks)
		wantErr   bool
		errMsg    string
		validate  func(t *testing.T, ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "successful create",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// No setup needed
			},
			msg: &types.MsgCreate{
				Creator:        testAddress,
				Admin:          testAddress,
				BaseYbtCreator: testAddress2,
			},
			setupMock: func(mocks moduleMocks) {
				// Mock base YBT token exists
				baseToken := ybtbasetypes.Token{
					Creator:      testAddress2,
					Admin:        testAddress3,
					Permissioned: false,
				}
				mocks.YbtbaseKeeper.EXPECT().
					GetToken(gomock.Any(), testAddress2).
					Return(baseToken, true)

				// Mock initial supply for collateral pool is zero
				collateralAddr := keeper.GetCollateralPoolAddress(testAddress)
				denom := keeper.GetTokenDenom(testAddress)
				mocks.BankKeeper.EXPECT().
					GetSupply(gomock.Any(), denom).
					Return(sdk.NewCoin(denom, math.ZeroInt()))

				// Mock set module account for collateral pool
				mocks.AuthKeeper.EXPECT().
					GetAccount(gomock.Any(), collateralAddr).
					Return(nil)
				mocks.AuthKeeper.EXPECT().
					SetModuleAccount(gomock.Any(), gomock.Any())
			},
			wantErr: false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check token was created
				token, found := k.GetToken(ctx, testAddress)
				require.True(t, found)
				require.Equal(t, testAddress, token.Creator)
				require.Equal(t, testAddress, token.Admin)
				require.Equal(t, testAddress2, token.BaseYbtCreator)
			},
		},
		{
			name: "token already exists",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create existing token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgCreate{
				Creator:        testAddress,
				Admin:          testAddress,
				BaseYbtCreator: testAddress2,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "token already exists",
		},
		{
			name: "base YBT token not found",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// No setup needed
			},
			msg: &types.MsgCreate{
				Creator:        testAddress,
				Admin:          testAddress,
				BaseYbtCreator: testAddress2,
			},
			setupMock: func(mocks moduleMocks) {
				// Mock base YBT token does not exist
				mocks.YbtbaseKeeper.EXPECT().
					GetToken(gomock.Any(), testAddress2).
					Return(ybtbasetypes.Token{}, false)
			},
			wantErr: true,
			errMsg:  "base YBT token not found",
		},
		{
			name: "invalid creator address",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// No setup needed
			},
			msg: &types.MsgCreate{
				Creator:        "invalid",
				Admin:          testAddress,
				BaseYbtCreator: testAddress2,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "invalid creator address",
		},
		{
			name: "invalid admin address",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// No setup needed
			},
			msg: &types.MsgCreate{
				Creator:        testAddress,
				Admin:          "invalid",
				BaseYbtCreator: testAddress2,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "invalid admin address",
		},
		{
			name: "invalid base YBT creator address",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// No setup needed
			},
			msg: &types.MsgCreate{
				Creator:        testAddress,
				Admin:          testAddress,
				BaseYbtCreator: "invalid",
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "invalid base YBT creator address",
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

			resp, err := ms.Create(ctx, tt.msg)
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