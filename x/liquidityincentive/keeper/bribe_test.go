package keeper

import (
	"strconv"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func setupKeeper(t *testing.T) (sdk.Context, Keeper) {
	ctx, k := setupKeeperWithParams(t)
	return ctx, k
}

func TestRegisterBribe(t *testing.T) {
	ctx, k := setupKeeper(t)
	msgServer := NewMsgServerImpl(k)

	// Create test accounts
	_, _, addr1 := testdata.KeyTestPubAddr()
	addr1Str := addr1.String()
	_, _, addr2 := testdata.KeyTestPubAddr()
	addr2Str := addr2.String()

	// Test cases
	tests := []struct {
		name      string
		msg       *types.MsgRegisterBribe
		expectErr bool
	}{
		{
			name: "valid bribe registration",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 1,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: false,
		},
		{
			name: "zero amount bribe",
			msg: &types.MsgRegisterBribe{
				Sender:  addr2Str,
				EpochId: 1,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(0))),
			},
			expectErr: true,
		},
		{
			name: "past epoch bribe",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 0,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expectErr {
				var senderAddr string
				if tc.msg.Sender == addr1Str {
					senderAddr = addr1Str
				} else if tc.msg.Sender == addr2Str {
					senderAddr = addr2Str
				} else {
					senderAddr = tc.msg.Sender
				}
				err := k.bankKeeper.(*MockBankKeeper).MintCoins(ctx, senderAddr, tc.msg.Amount)
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
	addr1Str := addr1.String()
	_, _, addr2 := testdata.KeyTestPubAddr()
	addr2Str := addr2.String()

	// Register a bribe
	bribeAmount := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))
	// Mint coins directly to addr1 before registering the bribe
	err := k.bankKeeper.(*MockBankKeeper).MintCoins(ctx, addr1Str, bribeAmount)
	require.NoError(t, err)

	msg := &types.MsgRegisterBribe{
		Sender:  addr1Str,
		EpochId: 1,
		PoolId:  1,
		Amount:  bribeAmount,
	}
	_, err = msgServer.RegisterBribe(ctx, msg)
	require.NoError(t, err)

	// Get bribe ID from event
	events := ctx.EventManager().Events()
	for i, event := range events {
		t.Logf("Event %d: Type=%s", i, event.Type)
		for j, attr := range event.Attributes {
			t.Logf("  Attr %d: Key=%s, Value=%s", j, attr.Key, string(attr.Value))
		}
	}
	var bribeId uint64
	for _, event := range events {
		if event.Type == "sunrise.liquidityincentive.v1.EventRegisterBribe" {
			for _, attr := range event.Attributes {
				if attr.Key == "id" {
					idStr, err := strconv.Unquote(string(attr.Value))
					if err != nil {
						idStr = string(attr.Value) // fallback if not quoted
					}
					parsedId, err := strconv.ParseUint(idStr, 10, 64)
					require.NoError(t, err)
					bribeId = parsedId
					break
				}
			}
		}
	}
	// require.NotZero(t, bribeId)

	// Create vote weights for bribe allocation
	vote := types.Vote{
		Sender: addr2.String(),
		PoolWeights: []types.PoolWeight{{
			PoolId: 1,
			Weight: "1.0",
		}},
	}
	err = k.SetVote(ctx, vote)
	require.NoError(t, err)

	// Save vote weights for bribes
	err = k.SaveVoteWeightsForBribes(ctx, 1)
	require.NoError(t, err)

	// Create epoch
	epoch := types.Epoch{
		Id:         1,
		StartBlock: ctx.BlockHeight(),
		EndBlock:   ctx.BlockHeight() + 100,
	}
	err = k.SetEpoch(ctx, epoch)
	require.NoError(t, err)

	// Create gauge for the pool
	gauge := types.Gauge{
		PoolId: 1,
		Count:  math.NewInt(1),
	}
	epoch.Gauges = append(epoch.Gauges, gauge)
	err = k.SetEpoch(ctx, epoch)
	require.NoError(t, err)

	// Create bribe allocation
	allocation := types.BribeAllocation{
		Address:         addr2Str,
		EpochId:         1,
		PoolId:          1,
		Weight:          "1.0",
		ClaimedBribeIds: []uint64{},
	}
	err = k.SetBribeAllocation(ctx, allocation)
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
				Sender:  addr2Str,
				BribeId: bribeId,
			},
			expectErr: false,
		},
		{
			name: "claim non-existent bribe",
			msg: &types.MsgClaimBribes{
				Sender:  addr2Str,
				BribeId: 999,
			},
			expectErr: true,
		},
		{
			name: "claim with wrong address",
			msg: &types.MsgClaimBribes{
				Sender:  addr1Str,
				BribeId: bribeId,
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
				allocation, err := k.GetBribeAllocation(ctx, addr2, 1, 1)
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
	addr1Str := addr1.String()

	// Register a bribe
	bribeAmount := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))
	// Mint coins directly to addr1 before registering the bribe
	err := k.bankKeeper.(*MockBankKeeper).MintCoins(ctx, addr1Str, bribeAmount)
	require.NoError(t, err)

	msg := &types.MsgRegisterBribe{
		Sender:  addr1Str,
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
