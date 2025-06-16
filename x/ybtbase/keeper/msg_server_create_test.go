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

func TestMsgServerCreate(t *testing.T) {
	tests := []struct {
		name      string
		msg       *types.MsgCreate
		setupMock func(mocks moduleMocks)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful creation - permissionless",
			msg: &types.MsgCreate{
				Creator:        testAddress,
				Admin:          testAddress,
				PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
			},
			setupMock: func(mocks moduleMocks) {
				// No bank operations needed for creation
			},
			wantErr: false,
		},
		{
			name: "successful creation - permissioned",
			msg: &types.MsgCreate{
				Creator:        testAddress,
				Admin:          testAddress2,
				PermissionMode: types.PermissionMode_PERMISSION_MODE_WHITELIST,
			},
			setupMock: func(mocks moduleMocks) {
				// SetSendEnabled should be called for non-permissionless tokens
				mocks.BankKeeper.EXPECT().SetSendEnabled(gomock.Any(), "ybtbase/"+testAddress, false)
			},
			wantErr: false,
		},
		{
			name: "token already exists",
			msg: &types.MsgCreate{
				Creator:        testAddress,
				Admin:          testAddress,
				PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
			},
			setupMock: func(mocks moduleMocks) {
				// No bank operations needed for creation
			},
			wantErr: false, // First creation should succeed, we'll test duplicate in the test body
		},
		{
			name: "invalid creator address",
			msg: &types.MsgCreate{
				Creator:        "invalid",
				Admin:          testAddress,
				PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "invalid creator address",
		},
		{
			name: "invalid admin address",
			msg: &types.MsgCreate{
				Creator:        testAddress,
				Admin:          "invalid",
				PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "invalid admin address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := initFixture(t)
			ctx := sdk.UnwrapSDKContext(f.ctx)
			k := f.keeper
			ms := keeper.NewMsgServerImpl(k)

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

				// Verify token was created
				token, found := k.GetToken(ctx, tt.msg.Creator)
				require.True(t, found)
				require.Equal(t, tt.msg.Creator, token.Creator)
				require.Equal(t, tt.msg.Admin, token.Admin)
				require.Equal(t, tt.msg.PermissionMode, token.PermissionMode)

				// Verify initial global reward index is set to 1
				globalIndex := k.GetGlobalRewardIndex(ctx, tt.msg.Creator)
				require.Equal(t, math.LegacyOneDec(), globalIndex)

				// Test duplicate creation
				if tt.name == "token already exists" {
					_, err := ms.Create(ctx, tt.msg)
					require.Error(t, err)
					require.Contains(t, err.Error(), "token already exists")
				}
			}
		})
	}
}
