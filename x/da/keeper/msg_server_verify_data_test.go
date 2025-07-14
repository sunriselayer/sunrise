package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestMsgServerVerifyData(t *testing.T) {
	k, mocks, srv, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Set up parameters
	params := types.DefaultParams()
	params.RejectedRemovalPeriod = 24 * time.Hour
	params.VerifiedRemovalPeriod = 24 * time.Hour
	params.ChallengePeriod = 24 * time.Hour
	params.ProofPeriod = 24 * time.Hour
	params.ChallengeThreshold = "0.5"
	params.ReplicationFactor = "2"
	params.SlashEpoch = 100
	require.NoError(t, k.Params.Set(sdkCtx, params))

	// Set up mocks for all potential calls within VerifyData
	mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(new(mockIterator), nil).AnyTimes()
	mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), gomock.Any()).Return(nil, stakingtypes.ErrNoValidatorFound).AnyTimes()
	mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Execute the message
	msg := &types.MsgVerifyData{
		Sender: sdk.AccAddress("sender").String(),
	}
	_, err := srv.VerifyData(ctx, msg)
	require.NoError(t, err)

	// Verify LastSlashBlockHeight was set
	lastSlashBlock, err := k.LastSlashBlockHeight.Get(sdkCtx)
	require.NoError(t, err)
	require.Equal(t, sdkCtx.BlockHeight(), lastSlashBlock)
}

func TestMsgServerVerifyData_SlashEpoch(t *testing.T) {
	k, mocks, srv, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Set up parameters
	params := types.DefaultParams()
	params.SlashEpoch = 50
	require.NoError(t, k.Params.Set(sdkCtx, params))

	// Set initial slash block height
	initialBlockHeight := sdkCtx.BlockHeight()
	require.NoError(t, k.LastSlashBlockHeight.Set(sdkCtx, initialBlockHeight))

	// Set up mocks
	mocks.StakingKeeper.EXPECT().TotalBondedTokens(gomock.Any()).Return(math.NewInt(1000), nil).AnyTimes()
	mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), gomock.Any()).Return(nil, stakingtypes.ErrNoValidatorFound).AnyTimes()
	mocks.StakingKeeper.EXPECT().PowerReduction(gomock.Any()).Return(sdk.DefaultPowerReduction).AnyTimes()
	mocks.StakingKeeper.EXPECT().ValidatorsPowerStoreIterator(gomock.Any()).Return(new(mockIterator), nil).AnyTimes()

	// Fast forward block height to trigger slash epoch
	sdkCtx = sdkCtx.WithBlockHeight(initialBlockHeight + int64(params.SlashEpoch))

	// Execute the message
	msg := &types.MsgVerifyData{
		Sender: sdk.AccAddress("sender").String(),
	}
	_, err := srv.VerifyData(sdk.WrapSDKContext(sdkCtx), msg)
	require.NoError(t, err)

	// Verify LastSlashBlockHeight was updated
	lastSlashBlock, err := k.LastSlashBlockHeight.Get(sdkCtx)
	require.NoError(t, err)
	require.Equal(t, sdkCtx.BlockHeight(), lastSlashBlock)
}
