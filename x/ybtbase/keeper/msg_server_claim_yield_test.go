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

func TestMsgServerClaimYield(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, k keeper.Keeper)
		msg       *types.MsgClaimYield
		setupMock func(mocks moduleMocks)
		wantErr   bool
		errMsg    string
		validate  func(t *testing.T, ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "successful claim yield - permissionless token",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create permissionless token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Set global reward index to 1.5 (50% yield accumulated)
				err = k.SetGlobalRewardIndex(ctx, testAddress, math.LegacyNewDecWithPrec(15, 1))
				require.NoError(t, err)

				// Set user's last reward index to 1.0
				err = k.SetUserLastRewardIndex(ctx, testAddress, testAddress2, math.LegacyOneDec())
				require.NoError(t, err)
			},
			msg: &types.MsgClaimYield{
				Sender:       testAddress2,
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {
				denom := keeper.GetTokenDenom(testAddress)
				userAddr := sdk.MustAccAddressFromBech32(testAddress2)
				yieldPoolAddr := keeper.GetYieldPoolAddress(testAddress)

				// Mock user balance - 1000 tokens
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), userAddr, denom).
					Return(sdk.NewCoin(denom, math.NewInt(1000)))

				// Expected yield: 1000 * (1.5 - 1.0) = 500
				yieldCoins := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(500)))

				// Mock yield pool balance check
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), yieldPoolAddr, denom).
					Return(sdk.NewCoin(denom, math.NewInt(2000)))

				// Expect transfer from yield pool to user
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), yieldPoolAddr, userAddr, yieldCoins).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check user's last reward index was updated
				index := k.GetUserLastRewardIndex(ctx, testAddress, testAddress2)
				require.Equal(t, math.LegacyNewDecWithPrec(15, 1), index)
			},
		},
		{
			name: "successful claim yield - permissioned token with permission",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create permissioned token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_WHITELIST,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Grant permission
				err = k.SetPermission(ctx, testAddress, testAddress2, true)
				require.NoError(t, err)

				// Set global reward index to 2.0 (100% yield accumulated)
				err = k.SetGlobalRewardIndex(ctx, testAddress, math.LegacyNewDec(2))
				require.NoError(t, err)

				// Set user's last reward index to 1.0
				err = k.SetUserLastRewardIndex(ctx, testAddress, testAddress2, math.LegacyOneDec())
				require.NoError(t, err)
			},
			msg: &types.MsgClaimYield{
				Sender:       testAddress2,
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {
				denom := keeper.GetTokenDenom(testAddress)
				userAddr := sdk.MustAccAddressFromBech32(testAddress2)
				yieldPoolAddr := keeper.GetYieldPoolAddress(testAddress)

				// Mock user balance - 500 tokens
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), userAddr, denom).
					Return(sdk.NewCoin(denom, math.NewInt(500)))

				// Expected yield: 500 * (2.0 - 1.0) = 500
				yieldCoins := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(500)))

				// Mock yield pool balance check
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), yieldPoolAddr, denom).
					Return(sdk.NewCoin(denom, math.NewInt(2000)))

				// Expect transfer from yield pool to user
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), yieldPoolAddr, userAddr, yieldCoins).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, ctx sdk.Context, k keeper.Keeper) {
				// Check user's last reward index was updated
				index := k.GetUserLastRewardIndex(ctx, testAddress, testAddress2)
				require.Equal(t, math.LegacyNewDec(2), index)
			},
		},
		{
			name: "no yield to claim - index already up to date",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Set both indexes to same value
				err = k.SetGlobalRewardIndex(ctx, testAddress, math.LegacyNewDec(2))
				require.NoError(t, err)
				err = k.SetUserLastRewardIndex(ctx, testAddress, testAddress2, math.LegacyNewDec(2))
				require.NoError(t, err)
			},
			msg: &types.MsgClaimYield{
				Sender:       testAddress2,
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {
				denom := keeper.GetTokenDenom(testAddress)
				userAddr := sdk.MustAccAddressFromBech32(testAddress2)

				// Mock user balance
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), userAddr, denom).
					Return(sdk.NewCoin(denom, math.NewInt(1000)))
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
					PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Set global reward index
				err = k.SetGlobalRewardIndex(ctx, testAddress, math.LegacyNewDec(2))
				require.NoError(t, err)
			},
			msg: &types.MsgClaimYield{
				Sender:       testAddress2,
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {
				denom := keeper.GetTokenDenom(testAddress)
				userAddr := sdk.MustAccAddressFromBech32(testAddress2)

				// Mock user balance - zero
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), userAddr, denom).
					Return(sdk.NewCoin(denom, math.ZeroInt()))
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
					PermissionMode: types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)

				// Set global reward index
				err = k.SetGlobalRewardIndex(ctx, testAddress, math.LegacyNewDec(2))
				require.NoError(t, err)

				// Set user's last reward index
				err = k.SetUserLastRewardIndex(ctx, testAddress, testAddress2, math.LegacyOneDec())
				require.NoError(t, err)
			},
			msg: &types.MsgClaimYield{
				Sender:       testAddress2,
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {
				denom := keeper.GetTokenDenom(testAddress)
				userAddr := sdk.MustAccAddressFromBech32(testAddress2)
				yieldPoolAddr := keeper.GetYieldPoolAddress(testAddress)

				// Mock user balance
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), userAddr, denom).
					Return(sdk.NewCoin(denom, math.NewInt(1000)))

				// Mock yield pool balance - insufficient
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), yieldPoolAddr, denom).
					Return(sdk.NewCoin(denom, math.NewInt(100)))
			},
			wantErr: true,
			errMsg:  "insufficient yield pool balance",
		},
		{
			name: "permissioned token without permission",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create permissioned token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					PermissionMode: types.PermissionMode_PERMISSION_MODE_WHITELIST,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
				// Don't grant permission
			},
			msg: &types.MsgClaimYield{
				Sender:       testAddress2,
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "no permission to claim yield",
		},
		{
			name: "token not found",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Don't create token
			},
			msg: &types.MsgClaimYield{
				Sender:       testAddress2,
				TokenCreator: testAddress,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "token not found",
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

			resp, err := ms.ClaimYield(ctx, tt.msg)
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
