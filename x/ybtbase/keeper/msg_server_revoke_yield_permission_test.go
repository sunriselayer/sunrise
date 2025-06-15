package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/ybtbase/keeper"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func TestMsgServerRevokeYieldPermission(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(ctx sdk.Context, k keeper.Keeper)
		msg      *types.MsgRevokeYieldPermission
		wantErr  bool
		errMsg   string
		validate func(t *testing.T, ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "successful revoke permission",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create permissioned token
				token := types.Token{
					Creator:      testAddress,
					Admin:        testAddress,
					Permissioned: true,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
				err = k.SetGlobalRewardIndex(ctx, testAddress, math.LegacyOneDec())
				require.NoError(t, err)
				// Grant permission first
				err = k.SetYieldPermission(ctx, testAddress, testAddress2, true)
				require.NoError(t, err)
			},
			msg: &types.MsgRevokeYieldPermission{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Target:       testAddress2,
			},
			wantErr: false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check permission was revoked
				hasPermission := k.HasYieldPermission(ctx, testAddress, testAddress2)
				require.False(t, hasPermission)
			},
		},
		{
			name: "revoke permission from non-permitted user",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create permissioned token
				token := types.Token{
					Creator:      testAddress,
					Admin:        testAddress,
					Permissioned: true,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
				// Don't grant permission
			},
			msg: &types.MsgRevokeYieldPermission{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Target:       testAddress2,
			},
			wantErr: false, // Should succeed even if not permitted
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check permission doesn't exist
				hasPermission := k.HasYieldPermission(ctx, testAddress, testAddress2)
				require.False(t, hasPermission)
			},
		},
		{
			name: "fail on permissionless token",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create permissionless token
				token := types.Token{
					Creator:      testAddress,
					Admin:        testAddress,
					Permissioned: false,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgRevokeYieldPermission{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Target:       testAddress2,
			},
			wantErr: true,
			errMsg:  "token is not permissioned",
		},
		{
			name: "token not found",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Don't create token
			},
			msg: &types.MsgRevokeYieldPermission{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Target:       testAddress2,
			},
			wantErr: true,
			errMsg:  "token not found",
		},
		{
			name: "unauthorized - not admin",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token with different admin
				token := types.Token{
					Creator:      testAddress,
					Admin:        testAddress2,
					Permissioned: true,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgRevokeYieldPermission{
				Admin:        testAddress, // Wrong admin
				TokenCreator: testAddress,
				Target:       testAddress3,
			},
			wantErr: true,
			errMsg:  "unauthorized",
		},
		{
			name: "invalid target address",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create permissioned token
				token := types.Token{
					Creator:      testAddress,
					Admin:        testAddress,
					Permissioned: true,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgRevokeYieldPermission{
				Admin:        testAddress,
				TokenCreator: testAddress,
				Target:       "invalid",
			},
			wantErr: true,
			errMsg:  "invalid target address",
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

			resp, err := ms.RevokeYieldPermission(ctx, tt.msg)
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