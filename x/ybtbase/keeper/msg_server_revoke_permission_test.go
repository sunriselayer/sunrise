package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/ybtbase/keeper"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func TestMsgRevokePermission(t *testing.T) {
	creator := testAddr1
	admin := testAddr2
	target := testAddr3
	nonAdmin := testAddr4

	tests := []struct {
		name      string
		msg       *types.MsgRevokePermission
		setupMock func(mocks moduleMocks)
		setupTest func(t *testing.T, k keeper.Keeper, ctx sdk.Context)
		wantErr   bool
		errMsg    string
		validate  func(t *testing.T, k keeper.Keeper, ctx sdk.Context)
	}{
		{
			name: "successful revoke permission - whitelist mode",
			msg: &types.MsgRevokePermission{
				Admin:        admin,
				TokenCreator: creator,
				Target:       target,
			},
			setupMock: func(mocks moduleMocks) {
				// When there are no more permissions, SetSendEnabled is called with false
				mocks.BankKeeper.EXPECT().SetSendEnabled(gomock.Any(), types.GetDenom(creator), false)
			},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// Create whitelist token
				token := types.Token{
					Creator:        creator,
					Admin:          admin,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_WHITELIST,
				}
				err := k.Tokens.Set(ctx, creator, token)
				require.NoError(t, err)
				
				// Set global reward index
				err = k.GlobalRewardIndex.Set(ctx, creator, math.LegacyOneDec())
				require.NoError(t, err)
				
				// Grant permission first
				err = k.Permissions.Set(ctx, collections.Join(creator, target), true)
				require.NoError(t, err)
			},
			wantErr: false,
			validate: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// Check permission was revoked (should not exist)
				hasPermission, err := k.Permissions.Has(ctx, collections.Join(creator, target))
				require.NoError(t, err)
				require.False(t, hasPermission)
			},
		},
		{
			name: "successful revoke permission - blacklist mode (add to blacklist)",
			msg: &types.MsgRevokePermission{
				Admin:        admin,
				TokenCreator: creator,
				Target:       target,
			},
			setupMock: func(mocks moduleMocks) {},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// Create blacklist token
				token := types.Token{
					Creator:        creator,
					Admin:          admin,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_BLACKLIST,
				}
				err := k.Tokens.Set(ctx, creator, token)
				require.NoError(t, err)
				
				// Set global reward index
				err = k.GlobalRewardIndex.Set(ctx, creator, math.LegacyOneDec())
				require.NoError(t, err)
			},
			wantErr: false,
			validate: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// Check target is blacklisted (should be true)
				hasPermission, err := k.Permissions.Has(ctx, collections.Join(creator, target))
				require.NoError(t, err)
				require.True(t, hasPermission)
				
				isBlocked, err := k.Permissions.Get(ctx, collections.Join(creator, target))
				require.NoError(t, err)
				require.True(t, isBlocked)
			},
		},
		{
			name: "fail - non-admin trying to revoke permission",
			msg: &types.MsgRevokePermission{
				Admin:        nonAdmin,
				TokenCreator: creator,
				Target:       target,
			},
			setupMock: func(mocks moduleMocks) {},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// Create whitelist token
				token := types.Token{
					Creator:        creator,
					Admin:          admin,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_WHITELIST,
				}
				err := k.Tokens.Set(ctx, creator, token)
				require.NoError(t, err)
			},
			wantErr: true,
			errMsg:  "unauthorized",
		},
		{
			name: "fail - token not found",
			msg: &types.MsgRevokePermission{
				Admin:        admin,
				TokenCreator: testAddr5, // non-existent token creator
				Target:       target,
			},
			setupMock: func(mocks moduleMocks) {},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {},
			wantErr:   true,
			errMsg:    "token not found",
		},
		{
			name: "fail - revoke permission on permissionless mode",
			msg: &types.MsgRevokePermission{
				Admin:        admin,
				TokenCreator: creator,
				Target:       target,
			},
			setupMock: func(mocks moduleMocks) {},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// Create permissionless token
				token := types.Token{
					Creator:        creator,
					Admin:          admin,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				err := k.Tokens.Set(ctx, creator, token)
				require.NoError(t, err)
			},
			wantErr: true,
			errMsg:  "cannot revoke permissions in permissionless mode",
		},
		{
			name: "fail - invalid admin address",
			msg: &types.MsgRevokePermission{
				Admin:        "invalid",
				TokenCreator: creator,
				Target:       target,
			},
			setupMock: func(mocks moduleMocks) {},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {},
			wantErr:   true,
			errMsg:    "invalid admin address",
		},
		{
			name: "fail - invalid token creator address",
			msg: &types.MsgRevokePermission{
				Admin:        admin,
				TokenCreator: "invalid",
				Target:       target,
			},
			setupMock: func(mocks moduleMocks) {},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {},
			wantErr:   true,
			errMsg:    "invalid token creator address",
		},
		{
			name: "fail - invalid target address",
			msg: &types.MsgRevokePermission{
				Admin:        admin,
				TokenCreator: creator,
				Target:       "invalid",
			},
			setupMock: func(mocks moduleMocks) {},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {},
			wantErr:   true,
			errMsg:    "invalid target address",
		},
		{
			name: "revoke non-existent permission in whitelist mode",
			msg: &types.MsgRevokePermission{
				Admin:        admin,
				TokenCreator: creator,
				Target:       target,
			},
			setupMock: func(mocks moduleMocks) {
				// SetSendEnabled is still called even if permission didn't exist
				mocks.BankKeeper.EXPECT().SetSendEnabled(gomock.Any(), types.GetDenom(creator), false)
			},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// Create whitelist token
				token := types.Token{
					Creator:        creator,
					Admin:          admin,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_WHITELIST,
				}
				err := k.Tokens.Set(ctx, creator, token)
				require.NoError(t, err)
				
				// Set global reward index
				err = k.GlobalRewardIndex.Set(ctx, creator, math.LegacyOneDec())
				require.NoError(t, err)
				
				// Don't grant any permission
			},
			wantErr: false, // Should succeed even if permission didn't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := initFixture(t)
			ctx := sdk.UnwrapSDKContext(f.ctx)
			k := f.keeper
			ms := keeper.NewMsgServerImpl(k)

			tt.setupMock(f.mocks)
			tt.setupTest(t, k, ctx)

			resp, err := ms.RevokePermission(ctx, tt.msg)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					require.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				
				if tt.validate != nil {
					tt.validate(t, k, ctx)
				}
			}
		})
	}
}