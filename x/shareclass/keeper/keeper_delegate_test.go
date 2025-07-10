package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
	tokenconvertertypes "github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func TestKeeper_Delegate(t *testing.T) {
	sender := sdk.AccAddress("sender")
	valAddr := sdk.ValAddress("validator")
	bondDenom := "bond"
	transferableDenom := "transferable"

	tests := []struct {
		name        string
		setup       func(f *fixture)
		sender      sdk.AccAddress
		valAddr     sdk.ValAddress
		amount      sdk.Coin
		expectedErr error
	}{
		{
			name: "success",
			setup: func(f *fixture) {
				amount := sdk.NewCoin(transferableDenom, math.NewInt(100))
				// Mock for ClaimRewards
				rewardSaverAddress := types.RewardSaverAddress(valAddr.String())
				f.mocks.BankKeeper.EXPECT().GetAllBalances(gomock.Any(), rewardSaverAddress).Return(sdk.NewCoins()).AnyTimes()

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

				moduleAddr := sdk.AccAddress("shareclass")
				f.mocks.AccountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr).AnyTimes()
				f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(amount)).Return(nil)
				f.mocks.TokenConverterKeeper.EXPECT().ConvertReverse(gomock.Any(), amount.Amount, moduleAddr).Return(nil)
				f.mocks.StakingMsgServer.EXPECT().Delegate(gomock.Any(), gomock.Any()).Return(&stakingtypes.MsgDelegateResponse{}, nil)

				// Mocks for Delegate
				f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil)
				f.mocks.BankKeeper.EXPECT().SetSendEnabled(gomock.Any(), types.NonVotingShareTokenDenom(valAddr.String()), false)
				f.mocks.BankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				f.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, sender, gomock.Any()).Return(nil)
			},
			sender:  sender,
			valAddr: valAddr,
			amount:  sdk.NewCoin(transferableDenom, math.NewInt(100)),
		},
		{
			name: "failure: invalid denom for delegate",
			setup: func(f *fixture) {
				// Mock GetTransferableDenom to return a different denom than the one in the test case
				f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return("expectedtransferable", nil)
			},
			sender:      sender,
			valAddr:     valAddr,
			amount:      sdk.NewCoin("wrongdenom", math.NewInt(100)), // This denom is different
			expectedErr: sdkerrors.ErrInvalidCoins,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.setup(f)

			_, _, err := f.keeper.Delegate(f.ctx, tc.sender, tc.valAddr, tc.amount)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_Undelegate(t *testing.T) {
	sender := sdk.AccAddress("sender")
	recipient := sdk.AccAddress("recipient")
	valAddr := sdk.ValAddress("validator")
	bondDenom := "bond"
	transferableDenom := "transferable"
	shareDenom := types.NonVotingShareTokenDenom(valAddr.String())

	tests := []struct {
		name        string
		setup       func(f *fixture)
		sender      sdk.AccAddress
		recipient   sdk.AccAddress
		valAddr     sdk.ValAddress
		amount      sdk.Coin
		expectedErr error
	}{
		{
			name: "success",
			setup: func(f *fixture) {
				amount := sdk.NewCoin(transferableDenom, math.NewInt(100))
				// Mock for ClaimRewards
				rewardSaverAddress := types.RewardSaverAddress(valAddr.String())
				f.mocks.BankKeeper.EXPECT().GetAllBalances(gomock.Any(), rewardSaverAddress).Return(sdk.NewCoins()).AnyTimes()

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
				f.mocks.AccountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.AccAddress("shareclass")).AnyTimes()
				f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)

				// Mocks for staking undelegate
				completionTime := time.Now()
				f.mocks.StakingMsgServer.EXPECT().Undelegate(gomock.Any(), gomock.Any()).Return(&stakingtypes.MsgUndelegateResponse{
					CompletionTime: completionTime,
					Amount:         amount,
				}, nil)
			},
			sender:    sender,
			recipient: recipient,
			valAddr:   valAddr,
			amount:    sdk.NewCoin(transferableDenom, math.NewInt(100)),
		},
		{
			name: "failure: invalid denom for undelegate",
			setup: func(f *fixture) {
				// Mock GetParams to return a different transferable denom
				f.mocks.TokenConverterKeeper.EXPECT().GetParams(gomock.Any()).Return(tokenconvertertypes.Params{
					NonTransferableDenom: bondDenom,
					TransferableDenom:    "expectedtransferable",
				}, nil)
				f.mocks.StakingKeeper.EXPECT().BondDenom(gomock.Any()).Return(bondDenom, nil)
			},
			sender:      sender,
			recipient:   recipient,
			valAddr:     valAddr,
			amount:      sdk.NewCoin("wrongdenom", math.NewInt(100)),
			expectedErr: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "failure: non-positive undelegate amount",
			setup: func(f *fixture) {
				// Mock GetParams to return correct transferable denom
				f.mocks.TokenConverterKeeper.EXPECT().GetParams(gomock.Any()).Return(tokenconvertertypes.Params{
					NonTransferableDenom: bondDenom,
					TransferableDenom:    transferableDenom,
				}, nil)
				f.mocks.StakingKeeper.EXPECT().BondDenom(gomock.Any()).Return(bondDenom, nil)
			},
			sender:      sender,
			recipient:   recipient,
			valAddr:     valAddr,
			amount:      sdk.NewCoin(transferableDenom, math.NewInt(0)),
			expectedErr: sdkerrors.ErrInvalidCoins,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.setup(f)

			_, _, _, err := f.keeper.Undelegate(f.ctx, tc.sender, tc.recipient, tc.valAddr, tc.amount)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
