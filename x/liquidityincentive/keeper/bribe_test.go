package keeper

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func TestRegisterBribe(t *testing.T) {
	ctx, k := setupKeeper(t)
	msgServer := NewMsgServerImpl(k)

	// Create test accounts
	_, _, addr1 := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()

	// Test cases
	tests := []struct {
		name      string
		msg       *types.MsgRegisterBribe
		expectErr bool
	}{
		{
			name: "valid bribe registration",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1.String(),
				EpochId: 1,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: false,
		},
		{
			name: "zero amount bribe",
			msg: &types.MsgRegisterBribe{
				Sender:  addr2.String(),
				EpochId: 1,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(0))),
			},
			expectErr: true,
		},
		{
			name: "past epoch bribe",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1.String(),
				EpochId: 0,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Fund the sender account with enough coins for valid case
			if tc.name == "valid bribe registration" {
				err := k.bankKeeper.(*MockBankKeeper).MintCoins(ctx, tc.msg.Sender, tc.msg.Amount)
				require.NoError(t, err)
			}
			// Fund the sender account for other cases if needed
			if !tc.expectErr && tc.name != "valid bribe registration" {
				err := k.bankKeeper.(*MockBankKeeper).MintCoins(ctx, types.ModuleName, tc.msg.Amount)
				require.NoError(t, err)
				err = k.bankKeeper.(*MockBankKeeper).SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(tc.msg.Sender), tc.msg.Amount)
				require.NoError(t, err)
			}

			// Register bribe
			_, err := msgServer.RegisterBribe(ctx, tc.msg)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify bribe was registered
				bribes, err := k.GetAllBribeByEpochId(ctx, tc.msg.EpochId)
				require.NoError(t, err)
				require.Len(t, bribes, 1)
				require.Equal(t, tc.msg.PoolId, bribes[0].PoolId)
				require.Equal(t, tc.msg.Amount, bribes[0].Amount)
				require.Equal(t, tc.msg.Sender, bribes[0].Address)
			}
		})
	}
}

func TestClaimBribes(t *testing.T) {
	ctx, k := setupKeeper(t)
	msgServer := NewMsgServerImpl(k)

	// Create test accounts
	_, _, addr1 := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()

	// Register a bribe
	bribeAmount := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))
	// Fund addr1 with enough coins before registering the bribe
	err := k.bankKeeper.(*MockBankKeeper).MintCoins(ctx, addr1.String(), bribeAmount)
	require.NoError(t, err)
	err = k.bankKeeper.(*MockBankKeeper).MintCoins(ctx, types.ModuleName, bribeAmount)
	require.NoError(t, err)
	err = k.bankKeeper.(*MockBankKeeper).SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr1, bribeAmount)
	require.NoError(t, err)

	msg := &types.MsgRegisterBribe{
		Sender:  addr1.String(),
		EpochId: 1,
		PoolId:  1,
		Amount:  bribeAmount,
	}
	_, err = msgServer.RegisterBribe(ctx, msg)
	require.NoError(t, err)

	// Create vote weights for bribe allocation
	vote := types.Vote{
		Sender: addr2.String(),
		PoolWeights: []types.PoolWeight{
			{
				PoolId: 1,
				Weight: "1.0",
			},
		},
	}
	err = k.SetVote(ctx, vote)
	require.NoError(t, err)

	// Save vote weights for bribes
	err = k.SaveVoteWeightsForBribes(ctx, 1)
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name      string
		msg       *types.MsgClaimBribes
		expectErr bool
	}{
		{
			name: "valid bribe claim",
			msg: &types.MsgClaimBribes{
				Sender:  addr2.String(),
				BribeId: 1,
			},
			expectErr: false,
		},
		{
			name: "claim non-existent bribe",
			msg: &types.MsgClaimBribes{
				Sender:  addr2.String(),
				BribeId: 999,
			},
			expectErr: true,
		},
		{
			name: "claim with wrong address",
			msg: &types.MsgClaimBribes{
				Sender:  addr1.String(),
				BribeId: 1,
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Claim bribe
			_, err := msgServer.ClaimBribes(ctx, tc.msg)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify bribe was claimed
				bribe, found, err := k.GetBribe(ctx, tc.msg.BribeId)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, bribeAmount, bribe.ClaimedAmount)

				// Verify allocation was updated
				allocation, err := k.GetBribeAllocation(ctx, sdk.AccAddress(tc.msg.Sender), 1, 1)
				require.NoError(t, err)
				require.Contains(t, allocation.ClaimedBribeIds, tc.msg.BribeId)
			}
		})
	}
}

func TestProcessUnclaimedBribes(t *testing.T) {
	ctx, k := setupKeeper(t)
	msgServer := NewMsgServerImpl(k)

	// Create test accounts
	_, _, addr1 := testdata.KeyTestPubAddr()

	// Register a bribe
	bribeAmount := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))
	err := k.bankKeeper.(*MockBankKeeper).MintCoins(ctx, types.ModuleName, bribeAmount)
	require.NoError(t, err)
	err = k.bankKeeper.(*MockBankKeeper).SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr1, bribeAmount)
	require.NoError(t, err)

	msg := &types.MsgRegisterBribe{
		Sender:  addr1.String(),
		EpochId: 1,
		PoolId:  1,
		Amount:  bribeAmount,
	}
	_, err = msgServer.RegisterBribe(ctx, msg)
	require.NoError(t, err)

	// Set expired epoch ID
	err = k.SetBribeExpiredEpochId(ctx, 0)
	require.NoError(t, err)

	// Process unclaimed bribes
	err = k.ProcessUnclaimedBribes(ctx, 1)
	require.NoError(t, err)

	// Verify bribe was removed
	bribes, err := k.GetAllBribeByEpochId(ctx, 1)
	require.NoError(t, err)
	require.Len(t, bribes, 0)

	// Verify expired epoch ID was updated
	expiredEpochId := k.GetBribeExpiredEpochId(ctx)
	require.Equal(t, uint64(1), expiredEpochId)
}

func setupKeeper(t *testing.T) (sdk.Context, Keeper) {
	// Create a new context and keeper for each test
	ctx, k := setupKeeperWithParams(t)
	return ctx, k
}
