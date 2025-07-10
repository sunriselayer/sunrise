package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
	"go.uber.org/mock/gomock"
)

func TestKeeper_GetShare(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	addr := sdk.AccAddress("address")
	valAddr := sdk.ValAddress("validator").String()
	shareDenom := types.NonVotingShareTokenDenom(valAddr)
	expectedShare := math.NewInt(100)

	f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), addr, shareDenom).Return(sdk.NewCoin(shareDenom, expectedShare))

	share := f.keeper.GetShare(f.ctx, addr, valAddr)
	require.Equal(expectedShare, share)
}

func TestKeeper_GetTotalShare(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	valAddr := sdk.ValAddress("validator").String()
	shareDenom := types.NonVotingShareTokenDenom(valAddr)
	expectedTotalShare := math.NewInt(1000)

	f.mocks.BankKeeper.EXPECT().GetSupply(gomock.Any(), shareDenom).Return(sdk.NewCoin(shareDenom, expectedTotalShare))

	totalShare := f.keeper.GetTotalShare(f.ctx, valAddr)
	require.Equal(expectedTotalShare, totalShare)
}

func TestKeeper_CalculateAmountByShare(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	valAddr := sdk.ValAddress("validator").String()
	share := math.NewInt(100)
	totalShare := math.NewInt(1000)
	totalStaked := math.NewInt(2000)

	shareDenom := types.NonVotingShareTokenDenom(valAddr)
	f.mocks.BankKeeper.EXPECT().GetSupply(gomock.Any(), shareDenom).Return(sdk.NewCoin(shareDenom, totalShare))

	moduleAddr := sdk.AccAddress("shareclass")
	f.mocks.AccountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr)
	f.mocks.StakingQueryServer.EXPECT().Delegation(gomock.Any(), &stakingtypes.QueryDelegationRequest{DelegatorAddr: moduleAddr.String(), ValidatorAddr: valAddr}).
		Return(&stakingtypes.QueryDelegationResponse{DelegationResponse: &stakingtypes.DelegationResponse{Balance: sdk.NewCoin("stake", totalStaked)}}, nil)

	amount, err := f.keeper.CalculateAmountByShare(f.ctx, valAddr, share)
	require.NoError(err)

	expectedAmount, err := types.CalculateAmountByShare(totalShare, totalStaked, share)
	require.NoError(err)
	require.Equal(expectedAmount, amount)
}

func TestKeeper_CalculateShareByAmount(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	valAddr := sdk.ValAddress("validator").String()
	amount := math.NewInt(100)
	totalShare := math.NewInt(1000)
	totalStaked := math.NewInt(2000)

	shareDenom := types.NonVotingShareTokenDenom(valAddr)
	f.mocks.BankKeeper.EXPECT().GetSupply(gomock.Any(), shareDenom).Return(sdk.NewCoin(shareDenom, totalShare))

	moduleAddr := sdk.AccAddress("shareclass")
	f.mocks.AccountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr)
	f.mocks.StakingQueryServer.EXPECT().Delegation(gomock.Any(), &stakingtypes.QueryDelegationRequest{DelegatorAddr: moduleAddr.String(), ValidatorAddr: valAddr}).
		Return(&stakingtypes.QueryDelegationResponse{DelegationResponse: &stakingtypes.DelegationResponse{Balance: sdk.NewCoin("stake", totalStaked)}}, nil)

	share, err := f.keeper.CalculateShareByAmount(f.ctx, valAddr, amount)
	require.NoError(err)

	expectedShare, err := types.CalculateShareByAmount(totalShare, totalStaked, amount)
	require.NoError(err)
	require.Equal(expectedShare, share)
}

func TestKeeper_CalculateShareByAmount_ZeroTotalShare(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	valAddr := sdk.ValAddress("validator").String()
	amount := math.NewInt(100)
	totalShare := math.NewInt(0)

	shareDenom := types.NonVotingShareTokenDenom(valAddr)
	f.mocks.BankKeeper.EXPECT().GetSupply(gomock.Any(), shareDenom).Return(sdk.NewCoin(shareDenom, totalShare))

	share, err := f.keeper.CalculateShareByAmount(f.ctx, valAddr, amount)
	require.NoError(err)
	require.Equal(amount, share)
}
