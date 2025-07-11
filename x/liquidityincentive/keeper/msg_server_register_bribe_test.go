package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	liquiditypooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestRegisterBribe(t *testing.T) {
	// Create test accounts
	_, _, addr1 := testdata.KeyTestPubAddr()
	addr1Str := addr1.String()
	_, _, addr2 := testdata.KeyTestPubAddr()
	addr2Str := addr2.String()

	bech32Codec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	// Test cases
	tests := []struct {
		name      string
		msg       *types.MsgRegisterBribe
		expectErr bool
		setup     func(fx *fixture, ctx sdk.Context)
	}{
		{
			name: "valid bribe registration",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 4, // Future epoch
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: false,
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				fx.mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), uint64(1)).Return(liquiditypooltypes.Pool{}, true, nil).AnyTimes()
				fx.mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil).AnyTimes()
				fx.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.BribeAccount, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil).AnyTimes()
				sdkCtx := ctx.WithBlockHeight(1000)
				err := setupEpochs(sdkCtx, &fx.keeper)
				require.NoError(t, err)
				err = fx.keeper.SetBribeCount(sdkCtx, 1)
				require.NoError(t, err)
			},
		},
		{
			name: "zero amount bribe",
			msg: &types.MsgRegisterBribe{
				Sender:  addr2Str,
				EpochId: 4,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(0))),
			},
			expectErr: true,
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				fx.mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), uint64(1)).Return(liquiditypooltypes.Pool{}, true, nil).AnyTimes()
				fx.mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil).AnyTimes()
				fx.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.BribeAccount, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil).AnyTimes()
				sdkCtx := ctx.WithBlockHeight(1000)
				err := setupEpochs(sdkCtx, &fx.keeper)
				require.NoError(t, err)
				err = fx.keeper.SetBribeCount(sdkCtx, 1)
				require.NoError(t, err)
			},
		},
		{
			name: "past epoch bribe",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 1,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: true,
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				sdkCtx := ctx.WithBlockHeight(1000)
				err := setupEpochs(sdkCtx, &fx.keeper)
				require.NoError(t, err)
				err = fx.keeper.SetBribeCount(sdkCtx, 1)
				require.NoError(t, err)
			},
		},
		{
			name: "current epoch bribe",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 2,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: true,
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				sdkCtx := ctx.WithBlockHeight(1000)
				err := setupEpochs(sdkCtx, &fx.keeper)
				require.NoError(t, err)
				err = fx.keeper.SetBribeCount(sdkCtx, 1)
				require.NoError(t, err)
			},
		},
		{
			name: "invalid pool id",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 4,
				PoolId:  999,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: true,
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				fx.mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), uint64(999)).Return(liquiditypooltypes.Pool{}, false, nil).AnyTimes()
				fx.mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil).AnyTimes()
				sdkCtx := ctx.WithBlockHeight(1000)
				err := setupEpochs(sdkCtx, &fx.keeper)
				require.NoError(t, err)
				err = fx.keeper.SetBribeCount(sdkCtx, 1)
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fx := initFixture(t)
			sdkCtx := fx.ctx.(sdk.Context)
			msgServer := keeper.NewMsgServerImpl(fx.keeper)

			// Set up test-specific mocks and keeper state
			if tc.setup != nil {
				tc.setup(fx, sdkCtx)
			}

			// Reset event manager for clean event tracking
			sdkCtx = sdkCtx.WithEventManager(sdk.NewEventManager())

			_, err := msgServer.RegisterBribe(sdkCtx, tc.msg)
			if tc.expectErr {
				require.Error(t, err)
				if tc.name == "zero amount bribe" {
					require.ErrorIs(t, err, types.ErrInvalidBribe)
					require.Contains(t, err.Error(), "amount cannot be zero")
				} else if tc.name == "invalid pool id" {
					require.ErrorIs(t, err, types.ErrInvalidBribe)
					require.Contains(t, err.Error(), "pool 999 not found")
				} else {
					require.ErrorIs(t, err, types.ErrInvalidBribe)
					require.Contains(t, err.Error(), "epoch must be in the future")
				}
			} else {
				require.NoError(t, err)
				bribes, err := fx.keeper.GetAllBribeByEpochId(sdkCtx, tc.msg.EpochId)
				require.NoError(t, err)
				require.Len(t, bribes, 1)
				require.Equal(t, tc.msg.PoolId, bribes[0].PoolId)
				require.Equal(t, tc.msg.Amount, bribes[0].Amount)
				require.Equal(t, tc.msg.Sender, bribes[0].Address)
				bribeId := bribes[0].Id
				require.Equal(t, uint64(1), bribeId) // Ensure BribeId starts at 1
			}
		})
	}
}
