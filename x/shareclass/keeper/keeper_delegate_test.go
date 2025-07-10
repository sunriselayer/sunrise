package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
	tokenconvertertypes "github.com/sunriselayer/sunrise/x/tokenconverter/types"
	"go.uber.org/mock/gomock"
)

func TestKeeper_Delegate(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	sender := sdk.AccAddress([]byte("sender"))
	valAddr := sdk.ValAddress([]byte("validator"))
	amount := sdk.NewCoin("transferable", math.NewInt(100))
	bondDenom := "bond"
	transferableDenom := "transferable"

	// Mock for ClaimRewards
	// It will call GetAllBalances, GetBalance, SendCoins
	rewardSaverAddress := types.RewardSaverAddress(valAddr.String())
	f.mocks.BankKeeper.EXPECT().GetAllBalances(gomock.Any(), rewardSaverAddress).Return(sdk.NewCoins(sdk.NewCoin(transferableDenom, math.NewInt(10)))).AnyTimes()
	f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoin(types.NonVotingShareTokenDenom(valAddr.String()), math.NewInt(1000))).AnyTimes()
	f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), rewardSaverAddress, sender, gomock.Any()).Return(nil)

	// Mock for GetTotalStakedAmount & CalculateShareByAmount
	f.mocks.StakingQueryServer.EXPECT().Delegation(gomock.Any(), gomock.Any()).Return(&stakingtypes.QueryDelegationResponse{
		DelegationResponse: &stakingtypes.DelegationResponse{
			Balance: sdk.NewCoin(bondDenom, math.NewInt(1000)),
		},
	}, nil).AnyTimes()
	f.mocks.BankKeeper.EXPECT().GetSupply(gomock.Any(), gomock.Any()).Return(sdk.NewCoin(types.NonVotingShareTokenDenom(valAddr.String()), math.NewInt(1000))).AnyTimes()

	// Mocks for ConvertAndDelegate
	f.mocks.StakingKeeper.EXPECT().BondDenom(gomock.Any()).Return(bondDenom, nil)
	f.mocks.TokenConverterKeeper.EXPECT().GetParams(gomock.Any()).Return(tokenconvertertypes.Params{
		NonTransferableDenom: bondDenom,
		TransferableDenom:    transferableDenom,
	}, nil)

	moduleAddr := sdk.AccAddress([]byte("module_address"))
	f.mocks.AccountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr).AnyTimes()
	f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(amount)).Return(nil)
	f.mocks.TokenConverterKeeper.EXPECT().ConvertReverse(gomock.Any(), amount.Amount, moduleAddr).Return(nil)
	f.mocks.StakingMsgServer.EXPECT().Delegate(gomock.Any(), gomock.Any()).Return(&stakingtypes.MsgDelegateResponse{}, nil)

	// Mocks for Delegate
	f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil)
	f.mocks.BankKeeper.EXPECT().SetSendEnabled(gomock.Any(), types.NonVotingShareTokenDenom(valAddr.String()), false)
	f.mocks.BankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
	f.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, sender, gomock.Any()).Return(nil)

	_, _, err := f.keeper.Delegate(f.ctx, sender, valAddr, amount)
	require.NoError(err)
}

func TestKeeper_Undelegate(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	sender := sdk.AccAddress([]byte("sender"))
	recipient := sdk.AccAddress([]byte("recipient"))
	valAddr := sdk.ValAddress([]byte("validator"))
	amount := sdk.NewCoin("transferable", math.NewInt(100))
	bondDenom := "bond"
	transferableDenom := "transferable"
	shareDenom := types.NonVotingShareTokenDenom(valAddr.String())

	// Mock for ClaimRewards
	rewardSaverAddress := types.RewardSaverAddress(valAddr.String())
	f.mocks.BankKeeper.EXPECT().GetAllBalances(gomock.Any(), rewardSaverAddress).Return(sdk.NewCoins(sdk.NewCoin(transferableDenom, math.NewInt(10)))).AnyTimes()
	f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoin(shareDenom, math.NewInt(1000)))
	f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), rewardSaverAddress, sender, gomock.Any()).Return(nil)

	// Mocks for Undelegate
	f.mocks.TokenConverterKeeper.EXPECT().GetParams(gomock.Any()).Return(tokenconvertertypes.Params{
		NonTransferableDenom: bondDenom,
		TransferableDenom:    transferableDenom,
	}, nil)
	f.mocks.StakingKeeper.EXPECT().BondDenom(gomock.Any()).Return(bondDenom, nil)

	// Mocks for CalculateShareByAmount
	f.mocks.StakingQueryServer.EXPECT().Delegation(gomock.Any(), gomock.Any()).Return(&stakingtypes.QueryDelegationResponse{
		DelegationResponse: &stakingtypes.DelegationResponse{
			Balance: sdk.NewCoin(bondDenom, math.NewInt(1000)),
		},
	}, nil).AnyTimes()
	f.mocks.BankKeeper.EXPECT().GetSupply(gomock.Any(), gomock.Any()).Return(sdk.NewCoin(shareDenom, math.NewInt(1000))).AnyTimes()

	// Mocks for bank send/burn
	f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, gomock.Any()).Return(nil)
	f.mocks.AccountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.AccAddress([]byte("module"))).AnyTimes()
	f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)

	// Mocks for staking undelegate
	completionTime := time.Now()
	f.mocks.StakingMsgServer.EXPECT().Undelegate(gomock.Any(), gomock.Any()).Return(&stakingtypes.MsgUndelegateResponse{
		CompletionTime: completionTime,
		Amount:         amount,
	}, nil)

	_, _, _, err := f.keeper.Undelegate(f.ctx, sender, recipient, valAddr, amount)
	require.NoError(err)
}
