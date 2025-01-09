package mint_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/header"
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/math"
	minttypes "cosmossdk.io/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/app/consts"
	"github.com/sunriselayer/sunrise/app/mint"
	liquidityincentivetypes "github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

type mockBankKeeper struct {
	mock.Mock
}

func (m *mockBankKeeper) GetSupply(ctx context.Context, denom string) sdk.Coin {
	args := m.Called(ctx, denom)
	return args.Get(0).(sdk.Coin)
}

func (m *mockBankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	args := m.Called(ctx, moduleName, amt)
	return args.Error(0)
}

func (m *mockBankKeeper) SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	args := m.Called(ctx, senderModule, recipientModule, amt)
	return args.Error(0)
}

type mockHeaderService struct {
	mock.Mock
}

func (m *mockHeaderService) HeaderInfo(ctx context.Context) header.Info {
	args := m.Called(ctx)
	return args.Get(0).(header.Info)
}

type mockQueryRouterService struct {
	mock.Mock
}

func (m *mockQueryRouterService) Invoke(ctx context.Context, req transaction.Msg) (transaction.Msg, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(transaction.Msg), args.Error(1)
}

func (m *mockQueryRouterService) CanInvoke(ctx context.Context, typeURL string) error {
	args := m.Called(ctx, typeURL)
	return args.Error(0)
}

func TestProvideMintFn(t *testing.T) {
	ctx := context.Background()
	mockBank := &mockBankKeeper{}
	mockHeader := &mockHeaderService{}
	mockRouter := &mockQueryRouterService{}

	env := appmodule.Environment{
		HeaderService:      mockHeader,
		QueryRouterService: mockRouter,
	}

	t.Run("skips non-minute epochs", func(t *testing.T) {
		mintFn := mint.ProvideMintFn(mockBank)
		minter := &minttypes.Minter{}

		err := mintFn(ctx, env, minter, "hour", 1)
		require.NoError(t, err)
	})

	t.Run("first time minting initialization", func(t *testing.T) {
		currentTime := time.Now()
		mockBank.On("GetSupply", ctx, consts.BondDenom).Return(sdk.NewCoin(consts.BondDenom, math.NewInt(1000000)))
		mockBank.On("GetSupply", ctx, consts.FeeDenom).Return(sdk.NewCoin(consts.FeeDenom, math.NewInt(0)))
		mockHeader.On("HeaderInfo", ctx).Return(header.Info{Time: currentTime})
		mockRouter.On("Invoke", ctx, mock.Anything).Return(&liquidityincentivetypes.QueryParamsResponse{
			Params: liquidityincentivetypes.Params{
				StakingRewardRatio: math.LegacyNewDecWithPrec(30, 2),
			},
		}, nil)
		mockBank.On("MintCoins", ctx, mock.Anything, mock.Anything).Return(nil)
		mockBank.On("SendCoinsFromModuleToModule", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		mintFn := mint.ProvideMintFn(mockBank)
		minter := &minttypes.Minter{}

		err := mintFn(ctx, env, minter, "minute", 1)
		require.NoError(t, err)
		require.NotNil(t, minter.Data)
		require.Equal(t, 8, len(minter.Data))
	})

	t.Run("zero provision when supply at cap", func(t *testing.T) {
		currentTime := time.Now()
		mockBank.On("GetSupply", ctx, consts.BondDenom).Return(sdk.NewCoin(consts.BondDenom, mint.SupplyCap))
		mockBank.On("GetSupply", ctx, consts.FeeDenom).Return(sdk.NewCoin(consts.FeeDenom, math.NewInt(0)))
		mockHeader.On("HeaderInfo", ctx).Return(header.Info{Time: currentTime})

		mintFn := mint.ProvideMintFn(mockBank)
		minter := &minttypes.Minter{Data: make([]byte, 8)}

		err := mintFn(ctx, env, minter, "minute", 1)
		require.NoError(t, err)
	})
}
