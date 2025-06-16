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

func TestMsgSend(t *testing.T) {
	creator := testAddr1
	admin := testAddr2
	fromAddr := testAddr3
	toAddr := testAddr4
	denom := types.GetDenom(creator)

	tests := []struct {
		name      string
		msg       *types.MsgSend
		setupMock func(mocks moduleMocks)
		setupTest func(t *testing.T, k keeper.Keeper, ctx sdk.Context)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful send - permissionless token",
			msg: &types.MsgSend{
				FromAddress:  fromAddr,
				TokenCreator: creator,
				ToAddress:    toAddr,
				Amount:       math.NewInt(500),
			},
			setupMock: func(mocks moduleMocks) {
				// Expect SendCoins to be called
				mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), 
					sdk.MustAccAddressFromBech32(fromAddr),
					sdk.MustAccAddressFromBech32(toAddr),
					sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(500)))).Return(nil)
			},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {
				// Create token
				token := types.Token{
					Creator:        creator,
					Admin:          admin,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				err := k.Tokens.Set(ctx, creator, token)
				require.NoError(t, err)
				
				// Set global reward index
				err = k.GlobalRewardIndex.Set(ctx, creator, math.LegacyOneDec())
				require.NoError(t, err)
			},
			wantErr: false,
		},
		{
			name: "successful send - whitelist mode with permissions",
			msg: &types.MsgSend{
				FromAddress:  admin,
				TokenCreator: creator,
				ToAddress:    fromAddr,
				Amount:       math.NewInt(300),
			},
			setupMock: func(mocks moduleMocks) {
				// Expect SendCoins to be called
				mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), 
					sdk.MustAccAddressFromBech32(admin),
					sdk.MustAccAddressFromBech32(fromAddr),
					sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(300)))).Return(nil)
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
				
				// Grant permissions (admin always has permission, grant to fromAddr)
				err = k.Permissions.Set(ctx, collections.Join(creator, fromAddr), true)
				require.NoError(t, err)
			},
			wantErr: false,
		},
		{
			name: "fail - whitelist mode without permission",
			msg: &types.MsgSend{
				FromAddress:  fromAddr,
				TokenCreator: creator,
				ToAddress:    toAddr,
				Amount:       math.NewInt(100),
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
				
				// Set global reward index
				err = k.GlobalRewardIndex.Set(ctx, creator, math.LegacyOneDec())
				require.NoError(t, err)
				
				// No permissions granted to fromAddr or toAddr
			},
			wantErr: true,
			errMsg:  "unauthorized",
		},
		{
			name: "fail - blacklist mode with blocked sender",
			msg: &types.MsgSend{
				FromAddress:  fromAddr,
				TokenCreator: creator,
				ToAddress:    toAddr,
				Amount:       math.NewInt(100),
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
				
				// Block fromAddr
				err = k.Permissions.Set(ctx, collections.Join(creator, fromAddr), true)
				require.NoError(t, err)
			},
			wantErr: true,
			errMsg:  "unauthorized",
		},
		{
			name: "fail - token not found",
			msg: &types.MsgSend{
				FromAddress:  fromAddr,
				TokenCreator: testAddr5, // non-existent token creator
				ToAddress:    toAddr,
				Amount:       math.NewInt(100),
			},
			setupMock: func(mocks moduleMocks) {},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {},
			wantErr:   true,
			errMsg:    "token not found",
		},
		{
			name: "fail - invalid from address",
			msg: &types.MsgSend{
				FromAddress:  "invalid",
				TokenCreator: creator,
				ToAddress:    toAddr,
				Amount:       math.NewInt(100),
			},
			setupMock: func(mocks moduleMocks) {},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {},
			wantErr:   true,
			errMsg:    "invalid from address",
		},
		{
			name: "fail - invalid recipient address",
			msg: &types.MsgSend{
				FromAddress:  fromAddr,
				TokenCreator: creator,
				ToAddress:    "invalid",
				Amount:       math.NewInt(100),
			},
			setupMock: func(mocks moduleMocks) {},
			setupTest: func(t *testing.T, k keeper.Keeper, ctx sdk.Context) {},
			wantErr:   true,
			errMsg:    "invalid to address",
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

			resp, err := ms.Send(ctx, tt.msg)
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