package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestChallengeCounter(t *testing.T) {
	k, _, _, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Initial value should be 0
	count := k.GetChallengeCounter(sdkCtx)
	require.Equal(t, uint64(0), count)

	// Set and get
	err := k.SetChallengeCounter(sdkCtx, 100)
	require.NoError(t, err)
	count = k.GetChallengeCounter(sdkCtx)
	require.Equal(t, uint64(100), count)
}

func TestFaultCounter(t *testing.T) {
	k, _, _, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr := sdk.ValAddress("validator")

	// Initial value should be 0
	count, err := k.GetFaultCounter(sdkCtx, valAddr)
	require.NoError(t, err)
	require.Equal(t, uint64(0), count)

	// Set and get
	err = k.SetFaultCounter(sdkCtx, valAddr, 5)
	require.NoError(t, err)
	count, err = k.GetFaultCounter(sdkCtx, valAddr)
	require.NoError(t, err)
	require.Equal(t, uint64(5), count)

	// Iterate
	var iteratedAddr sdk.ValAddress
	var iteratedCount uint64
	err = k.IterateFaultCounters(sdkCtx, func(operator sdk.ValAddress, faultCount uint64) (stop bool, err error) {
		iteratedAddr = operator
		iteratedCount = faultCount
		return true, nil
	})
	require.NoError(t, err)
	require.Equal(t, valAddr.String(), iteratedAddr.String())
	require.Equal(t, uint64(5), iteratedCount)

	// Delete
	err = k.DeleteFaultCounter(sdkCtx, valAddr)
	require.NoError(t, err)
	count, err = k.GetFaultCounter(sdkCtx, valAddr)
	require.NoError(t, err)
	require.Equal(t, uint64(0), count)
}

func TestHandleSlashEpoch(t *testing.T) {
	k, mocks, _, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Setup params and counters
	params := types.DefaultParams()
	params.SlashFaultThreshold = "0.5"
	params.SlashFraction = "0.01"
	require.NoError(t, k.Params.Set(sdkCtx, params))
	require.NoError(t, k.SetChallengeCounter(sdkCtx, 10))

	valAddr1 := sdk.ValAddress("validator1")
	valAddr2 := sdk.ValAddress("validator2")
	valAddr3 := sdk.ValAddress("validator3")

	// fault count > threshold
	require.NoError(t, k.SetFaultCounter(sdkCtx, valAddr1, 6))
	// fault count <= threshold
	require.NoError(t, k.SetFaultCounter(sdkCtx, valAddr2, 5))
	// validator not found
	require.NoError(t, k.SetFaultCounter(sdkCtx, valAddr3, 7))

	// Mock validator
	pubKey := ed25519.GenPrivKey().PubKey()
	anyPk, err := codectypes.NewAnyWithValue(pubKey)
	require.NoError(t, err)
	mockVal1 := &stakingtypes.Validator{
		OperatorAddress:   valAddr1.String(),
		ConsensusPubkey:   anyPk,
		Status:            stakingtypes.Bonded,
		Tokens:            math.NewInt(100),
		DelegatorShares:   math.LegacyNewDec(100),
		MinSelfDelegation: math.NewInt(1),
	}

	// Setup mocks
	mocks.StakingKeeper.EXPECT().PowerReduction(gomock.Any()).Return(sdk.DefaultPowerReduction).AnyTimes()
	mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr1).Return(mockVal1, nil).AnyTimes()
	mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr2).Return(mockVal1, nil).AnyTimes() // Re-use mockVal for simplicity
	mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr3).Return(nil, stakingtypes.ErrNoValidatorFound).AnyTimes()

	mocks.SlashingKeeper.EXPECT().Slash(gomock.Any(), pubKey.Address().Bytes(), math.LegacyMustNewDecFromStr(params.SlashFraction), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mocks.SlashingKeeper.EXPECT().Jail(gomock.Any(), pubKey.Address().Bytes()).Return(nil).Times(1)

	// Execute
	err = k.HandleSlashEpoch(sdkCtx)
	require.NoError(t, err)

	// Verify counters were reset/deleted
	challengeCount := k.GetChallengeCounter(sdkCtx)
	require.Equal(t, uint64(0), challengeCount)

	faultCount, err := k.GetFaultCounter(sdkCtx, valAddr1)
	require.NoError(t, err)
	require.Equal(t, uint64(0), faultCount)
}

func TestHandleSlashEpoch_JailedValidator(t *testing.T) {
	k, mocks, _, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Setup params and counters
	params := types.DefaultParams()
	params.SlashFaultThreshold = "0.5"
	params.SlashFraction = "0.01"
	require.NoError(t, k.Params.Set(sdkCtx, params))
	require.NoError(t, k.SetChallengeCounter(sdkCtx, 10))

	valAddr1 := sdk.ValAddress("validator1")

	// fault count > threshold
	require.NoError(t, k.SetFaultCounter(sdkCtx, valAddr1, 6))

	// Mock validator
	pubKey := ed25519.GenPrivKey().PubKey()
	anyPk, err := codectypes.NewAnyWithValue(pubKey)
	require.NoError(t, err)
	mockVal1 := &stakingtypes.Validator{
		OperatorAddress:   valAddr1.String(),
		ConsensusPubkey:   anyPk,
		Status:            stakingtypes.Bonded,
		Tokens:            math.NewInt(100),
		DelegatorShares:   math.LegacyNewDec(100),
		MinSelfDelegation: math.NewInt(1),
		Jailed:            true,
	}

	// Setup mocks
	mocks.StakingKeeper.EXPECT().PowerReduction(gomock.Any()).Return(sdk.DefaultPowerReduction).AnyTimes()
	mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr1).Return(mockVal1, nil).AnyTimes()

	// Execute
	err = k.HandleSlashEpoch(sdkCtx)
	require.NoError(t, err)

	// Verify counters were reset/deleted
	challengeCount := k.GetChallengeCounter(sdkCtx)
	require.Equal(t, uint64(0), challengeCount)

	faultCount, err := k.GetFaultCounter(sdkCtx, valAddr1)
	require.NoError(t, err)
	require.Equal(t, uint64(0), faultCount)
}

func TestHandleSlashEpoch_UnbondedValidator(t *testing.T) {
	k, mocks, _, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Setup params and counters
	params := types.DefaultParams()
	params.SlashFaultThreshold = "0.5"
	params.SlashFraction = "0.01"
	require.NoError(t, k.Params.Set(sdkCtx, params))
	require.NoError(t, k.SetChallengeCounter(sdkCtx, 10))

	valAddr1 := sdk.ValAddress("validator1")

	// fault count > threshold
	require.NoError(t, k.SetFaultCounter(sdkCtx, valAddr1, 6))

	// Mock validator
	pubKey := ed25519.GenPrivKey().PubKey()
	anyPk, err := codectypes.NewAnyWithValue(pubKey)
	require.NoError(t, err)
	mockVal1 := &stakingtypes.Validator{
		OperatorAddress:   valAddr1.String(),
		ConsensusPubkey:   anyPk,
		Status:            stakingtypes.Unbonded,
		Tokens:            math.NewInt(100),
		DelegatorShares:   math.LegacyNewDec(100),
		MinSelfDelegation: math.NewInt(1),
	}

	// Setup mocks
	mocks.StakingKeeper.EXPECT().PowerReduction(gomock.Any()).Return(sdk.DefaultPowerReduction).AnyTimes()
	mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr1).Return(mockVal1, nil).AnyTimes()

	// Execute
	err = k.HandleSlashEpoch(sdkCtx)
	require.NoError(t, err)

	// Verify counters were reset/deleted
	challengeCount := k.GetChallengeCounter(sdkCtx)
	require.Equal(t, uint64(0), challengeCount)

	faultCount, err := k.GetFaultCounter(sdkCtx, valAddr1)
	require.NoError(t, err)
	require.Equal(t, uint64(0), faultCount)
}

func TestHandleSlashEpoch_SlashError(t *testing.T) {
	k, mocks, _, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Setup params and counters
	params := types.DefaultParams()
	params.SlashFaultThreshold = "0.5"
	params.SlashFraction = "0.01"
	require.NoError(t, k.Params.Set(sdkCtx, params))
	require.NoError(t, k.SetChallengeCounter(sdkCtx, 10))

	valAddr1 := sdk.ValAddress("validator1")

	// fault count > threshold
	require.NoError(t, k.SetFaultCounter(sdkCtx, valAddr1, 6))

	// Mock validator
	pubKey := ed25519.GenPrivKey().PubKey()
	anyPk, err := codectypes.NewAnyWithValue(pubKey)
	require.NoError(t, err)
	mockVal1 := &stakingtypes.Validator{
		OperatorAddress:   valAddr1.String(),
		ConsensusPubkey:   anyPk,
		Status:            stakingtypes.Bonded,
		Tokens:            math.NewInt(100),
		DelegatorShares:   math.LegacyNewDec(100),
		MinSelfDelegation: math.NewInt(1),
	}

	// Setup mocks
	mocks.StakingKeeper.EXPECT().PowerReduction(gomock.Any()).Return(sdk.DefaultPowerReduction).AnyTimes()
	mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr1).Return(mockVal1, nil).AnyTimes()

	slashErr := errors.New("slash error")
	mocks.SlashingKeeper.EXPECT().Slash(gomock.Any(), pubKey.Address().Bytes(), math.LegacyMustNewDecFromStr(params.SlashFraction), gomock.Any(), gomock.Any()).Return(slashErr).Times(1)

	// Execute
	err = k.HandleSlashEpoch(sdkCtx)
	require.Error(t, err)
	require.Equal(t, slashErr, err)
}

func TestHandleSlashEpoch_JailError(t *testing.T) {
	k, mocks, _, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Setup params and counters
	params := types.DefaultParams()
	params.SlashFaultThreshold = "0.5"
	params.SlashFraction = "0.01"
	require.NoError(t, k.Params.Set(sdkCtx, params))
	require.NoError(t, k.SetChallengeCounter(sdkCtx, 10))

	valAddr1 := sdk.ValAddress("validator1")

	// fault count > threshold
	require.NoError(t, k.SetFaultCounter(sdkCtx, valAddr1, 6))

	// Mock validator
	pubKey := ed25519.GenPrivKey().PubKey()
	anyPk, err := codectypes.NewAnyWithValue(pubKey)
	require.NoError(t, err)
	mockVal1 := &stakingtypes.Validator{
		OperatorAddress:   valAddr1.String(),
		ConsensusPubkey:   anyPk,
		Status:            stakingtypes.Bonded,
		Tokens:            math.NewInt(100),
		DelegatorShares:   math.LegacyNewDec(100),
		MinSelfDelegation: math.NewInt(1),
	}

	// Setup mocks
	mocks.StakingKeeper.EXPECT().PowerReduction(gomock.Any()).Return(sdk.DefaultPowerReduction).AnyTimes()
	mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), valAddr1).Return(mockVal1, nil).AnyTimes()

	slashFraction, err := math.LegacyNewDecFromStr(params.SlashFraction)
	require.NoError(t, err)
	jailErr := errors.New("jail error")
	mocks.SlashingKeeper.EXPECT().Slash(gomock.Any(), pubKey.Address().Bytes(), slashFraction, gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mocks.SlashingKeeper.EXPECT().Jail(gomock.Any(), pubKey.Address().Bytes()).Return(jailErr).Times(1)

	// Execute
	err = k.HandleSlashEpoch(sdkCtx)
	require.Error(t, err)
	require.Equal(t, jailErr, err)
}

func TestHandleSlashEpoch_InvalidParams(t *testing.T) {
	testCases := []struct {
		name        string
		setupParams func(params *types.Params)
	}{
		{
			name: "invalid slash fault threshold",
			setupParams: func(params *types.Params) {
				params.SlashFaultThreshold = "invalid"
			},
		},
		{
			name: "invalid slash fraction",
			setupParams: func(params *types.Params) {
				params.SlashFraction = "invalid"
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k, _, _, ctx := setupMsgServer(t)
			sdkCtx := sdk.UnwrapSDKContext(ctx)

			params := types.DefaultParams()
			tc.setupParams(&params)
			require.NoError(t, k.Params.Set(sdkCtx, params))

			err := k.HandleSlashEpoch(sdkCtx)
			require.Error(t, err)
		})
	}
}

