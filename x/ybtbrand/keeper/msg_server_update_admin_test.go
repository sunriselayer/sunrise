package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/ybtbrand/keeper"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func TestMsgServerUpdateAdmin(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, k keeper.Keeper)
		msg       *types.MsgUpdateAdmin
		setupMock func(mocks moduleMocks)
		wantErr   bool
		errMsg    string
		validate  func(t *testing.T, ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "successful update admin",
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
			msg: &types.MsgUpdateAdmin{
				Admin:        testAddress,
				NewAdmin:     testAddress3,
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check admin was updated
				token, found := k.GetToken(ctx, testAddress)
				require.True(t, found)
				require.Equal(t, testAddress3, token.Admin)
			},
		},
		{
			name: "update admin to same address",
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
			msg: &types.MsgUpdateAdmin{
				Admin:        testAddress,
				NewAdmin:     testAddress,
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   false, // Should succeed even if same address
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check admin is still the same
				token, found := k.GetToken(ctx, testAddress)
				require.True(t, found)
				require.Equal(t, testAddress, token.Admin)
			},
		},
		{
			name: "token not found",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Don't create token
			},
			msg: &types.MsgUpdateAdmin{
				Admin:        testAddress,
				NewAdmin:     testAddress2,
				TokenCreator: testAddress,
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
			msg: &types.MsgUpdateAdmin{
				Admin:        testAddress, // Wrong admin
				NewAdmin:     testAddress3,
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "unauthorized",
		},
		{
			name: "invalid new admin address",
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
			msg: &types.MsgUpdateAdmin{
				Admin:        testAddress,
				NewAdmin:     "invalid",
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "invalid new admin address",
		},
		{
			name: "invalid admin address",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// No setup needed
			},
			msg: &types.MsgUpdateAdmin{
				Admin:        "invalid",
				NewAdmin:     testAddress2,
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "invalid admin address",
		},
		{
			name: "invalid token creator address",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// No setup needed
			},
			msg: &types.MsgUpdateAdmin{
				Admin:        testAddress,
				NewAdmin:     testAddress2,
				TokenCreator: "invalid",
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "invalid token creator address",
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

			resp, err := ms.UpdateAdmin(ctx, tt.msg)
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