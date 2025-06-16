package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	ybtbasetypes "github.com/sunriselayer/sunrise/x/ybtbase/types"
	"github.com/sunriselayer/sunrise/x/ybtbrand/keeper"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func TestMsgServerClaimCollateralYield(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx sdk.Context, k keeper.Keeper)
		msg       *types.MsgClaimCollateralYield
		setupMock func(mocks moduleMocks)
		wantErr   bool
		errMsg    string
		validate  func(t *testing.T, ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "successful claim - permissionless base YBT",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create brand token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgClaimCollateralYield{
				Admin:          testAddress,
				TokenCreator:   testAddress,
				BaseYbtCreator: testAddress2,
			},
			setupMock: func(mocks moduleMocks) {
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)
				collateralAddr := keeper.GetCollateralPoolAddress(testAddress)
				baseDenom := types.GetBaseYbtTokenDenom(testAddress2)
				
				// Mock get base token (permissionless)
				baseToken := ybtbasetypes.Token{
					Creator:        testAddress2,
					Admin:          testAddress3,
					PermissionMode: ybtbasetypes.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				mocks.YbtbaseKeeper.EXPECT().
					GetToken(gomock.Any(), testAddress2).
					Return(baseToken, true)
				
				// Mock collateral balance - 10000 base YBT
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), collateralAddr, baseDenom).
					Return(sdk.NewCoin(baseDenom, math.NewInt(10000)))
				
				// Mock get global reward index = 1.5
				mocks.YbtbaseKeeper.EXPECT().
					GetGlobalRewardIndex(gomock.Any(), testAddress2).
					Return(math.LegacyNewDecWithPrec(15, 1))
				
				// Mock get collateral pool's last reward index = 1.0
				mocks.YbtbaseKeeper.EXPECT().
					GetUserLastRewardIndex(gomock.Any(), testAddress2, collateralAddr.String()).
					Return(math.LegacyOneDec())
				
				// Expected yield: 10000 * (1.5 - 1.0) = 5000
				// Mock transfer yield from base YBT yield pool to admin
				yieldPoolAddr := types.GetBaseYbtYieldPoolAddress(testAddress2)
				yieldCoins := sdk.NewCoins(sdk.NewCoin(baseDenom, math.NewInt(5000)))
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), yieldPoolAddr, adminAddr, yieldCoins).
					Return(nil)
				
				// Mock update collateral pool's last reward index to 1.5
				mocks.YbtbaseKeeper.EXPECT().
					SetUserLastRewardIndex(gomock.Any(), testAddress2, collateralAddr.String(), math.LegacyNewDecWithPrec(15, 1)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "successful claim - permissioned base YBT with permission",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create brand token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgClaimCollateralYield{
				Admin:          testAddress,
				TokenCreator:   testAddress,
				BaseYbtCreator: testAddress2,
			},
			setupMock: func(mocks moduleMocks) {
				adminAddr := sdk.MustAccAddressFromBech32(testAddress)
				collateralAddr := keeper.GetCollateralPoolAddress(testAddress)
				baseDenom := "ybtbase/" + testAddress2
				
				// Mock get base token (permissioned)
				baseToken := ybtbasetypes.Token{
					Creator:      testAddress2,
					Admin:        testAddress3,
					PermissionMode: ybtbasetypes.PermissionMode_PERMISSION_MODE_WHITELIST,
				}
				mocks.YbtbaseKeeper.EXPECT().
					GetToken(gomock.Any(), testAddress2).
					Return(baseToken, true)
				
				// Mock check permission - admin has permission
				mocks.YbtbaseKeeper.EXPECT().
					HasPermission(gomock.Any(), testAddress2, testAddress).
					Return(true)
				
				// Mock collateral balance - 5000 base YBT
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), collateralAddr, baseDenom).
					Return(sdk.NewCoin(baseDenom, math.NewInt(5000)))
				
				// Mock get global reward index = 2.0
				mocks.YbtbaseKeeper.EXPECT().
					GetGlobalRewardIndex(gomock.Any(), testAddress2).
					Return(math.LegacyNewDec(2))
				
				// Mock get collateral pool's last reward index = 1.2
				mocks.YbtbaseKeeper.EXPECT().
					GetUserLastRewardIndex(gomock.Any(), testAddress2, collateralAddr.String()).
					Return(math.LegacyNewDecWithPrec(12, 1))
				
				// Expected yield: 5000 * (2.0 - 1.2) = 4000
				// Mock transfer yield from base YBT yield pool to admin
				yieldPoolAddr := types.GetBaseYbtYieldPoolAddress(testAddress2)
				yieldCoins := sdk.NewCoins(sdk.NewCoin(baseDenom, math.NewInt(4000)))
				mocks.BankKeeper.EXPECT().
					SendCoins(gomock.Any(), yieldPoolAddr, adminAddr, yieldCoins).
					Return(nil)
				
				// Mock update collateral pool's last reward index to 2.0
				mocks.YbtbaseKeeper.EXPECT().
					SetUserLastRewardIndex(gomock.Any(), testAddress2, collateralAddr.String(), math.LegacyNewDec(2)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "token not found",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Don't create token
			},
			msg: &types.MsgClaimCollateralYield{
				Admin:          testAddress,
				TokenCreator:   testAddress,
				BaseYbtCreator: testAddress2,
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
			msg: &types.MsgClaimCollateralYield{
				Admin:          testAddress, // Wrong admin
				TokenCreator:   testAddress,
				BaseYbtCreator: testAddress3,
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "unauthorized",
		},
		{
			name: "base YBT creator mismatch",
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
			msg: &types.MsgClaimCollateralYield{
				Admin:          testAddress,
				TokenCreator:   testAddress,
				BaseYbtCreator: testAddress3, // Wrong base YBT creator
			},
			setupMock: func(mocks moduleMocks) {},
			wantErr:   true,
			errMsg:    "base YBT creator mismatch",
		},
		{
			name: "base YBT token not found",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create brand token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgClaimCollateralYield{
				Admin:          testAddress,
				TokenCreator:   testAddress,
				BaseYbtCreator: testAddress2,
			},
			setupMock: func(mocks moduleMocks) {
				// Mock base token not found
				mocks.YbtbaseKeeper.EXPECT().
					GetToken(gomock.Any(), testAddress2).
					Return(ybtbasetypes.Token{}, false)
			},
			wantErr: true,
			errMsg:  "base YBT token not found",
		},
		{
			name: "no permission for permissioned base YBT",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create brand token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgClaimCollateralYield{
				Admin:          testAddress,
				TokenCreator:   testAddress,
				BaseYbtCreator: testAddress2,
			},
			setupMock: func(mocks moduleMocks) {
				// Mock get base token (permissioned)
				baseToken := ybtbasetypes.Token{
					Creator:        testAddress2,
					Admin:          testAddress3,
					PermissionMode: ybtbasetypes.PermissionMode_PERMISSION_MODE_WHITELIST,
				}
				mocks.YbtbaseKeeper.EXPECT().
					GetToken(gomock.Any(), testAddress2).
					Return(baseToken, true)
				
				// Mock check permission - admin has NO permission
				mocks.YbtbaseKeeper.EXPECT().
					HasPermission(gomock.Any(), testAddress2, testAddress).
					Return(false)
			},
			wantErr: true,
			errMsg:  "no permission for base YBT",
		},
		{
			name: "no collateral balance",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create brand token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgClaimCollateralYield{
				Admin:          testAddress,
				TokenCreator:   testAddress,
				BaseYbtCreator: testAddress2,
			},
			setupMock: func(mocks moduleMocks) {
				collateralAddr := keeper.GetCollateralPoolAddress(testAddress)
				baseDenom := "ybtbase/" + testAddress2
				
				// Mock get base token (permissionless)
				baseToken := ybtbasetypes.Token{
					Creator:        testAddress2,
					Admin:          testAddress3,
					PermissionMode: ybtbasetypes.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				mocks.YbtbaseKeeper.EXPECT().
					GetToken(gomock.Any(), testAddress2).
					Return(baseToken, true)
				
				// Mock zero collateral balance
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), collateralAddr, baseDenom).
					Return(sdk.NewCoin(baseDenom, math.ZeroInt()))
			},
			wantErr: true,
			errMsg:  "no collateral balance",
		},
		{
			name: "no yield to claim",
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				// Create brand token
				token := types.Token{
					Creator:        testAddress,
					Admin:          testAddress,
					BaseYbtCreator: testAddress2,
				}
				err := k.SetToken(ctx, testAddress, token)
				require.NoError(t, err)
			},
			msg: &types.MsgClaimCollateralYield{
				Admin:          testAddress,
				TokenCreator:   testAddress,
				BaseYbtCreator: testAddress2,
			},
			setupMock: func(mocks moduleMocks) {
				collateralAddr := keeper.GetCollateralPoolAddress(testAddress)
				baseDenom := "ybtbase/" + testAddress2
				
				// Mock get base token (permissionless)
				baseToken := ybtbasetypes.Token{
					Creator:        testAddress2,
					Admin:          testAddress3,
					PermissionMode: ybtbasetypes.PermissionMode_PERMISSION_MODE_PERMISSIONLESS,
				}
				mocks.YbtbaseKeeper.EXPECT().
					GetToken(gomock.Any(), testAddress2).
					Return(baseToken, true)
				
				// Mock collateral balance
				mocks.BankKeeper.EXPECT().
					GetBalance(gomock.Any(), collateralAddr, baseDenom).
					Return(sdk.NewCoin(baseDenom, math.NewInt(10000)))
				
				// Mock get global reward index = 1.0
				mocks.YbtbaseKeeper.EXPECT().
					GetGlobalRewardIndex(gomock.Any(), testAddress2).
					Return(math.LegacyOneDec())
				
				// Mock get collateral pool's last reward index = 1.0 (same as global)
				mocks.YbtbaseKeeper.EXPECT().
					GetUserLastRewardIndex(gomock.Any(), testAddress2, collateralAddr.String()).
					Return(math.LegacyOneDec())
				
				// Expected yield: 10000 * (1.0 - 1.0) = 0
			},
			wantErr: true,
			errMsg:  "no yield to claim",
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

			resp, err := ms.ClaimCollateralYield(ctx, tt.msg)
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