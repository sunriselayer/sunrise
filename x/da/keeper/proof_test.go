package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestProofStore(t *testing.T) {
	k, ctx := keepertest.DaKeeper(t)
	sender1 := sdk.AccAddress("sender1")
	sender2 := sdk.AccAddress("sender2")

	proof := k.GetProof(ctx, "ipfs://metadata1", sender1.String())
	require.Equal(t, proof.MetadataUri, "")

	err := k.SetProof(ctx, types.Proof{
		MetadataUri: "ipfs://metadata1",
		Sender:      sender1.String(),
		Indices:     []int64{0},
		Proofs:      [][]byte{{0x0}},
		IsValidData: true,
	})
	require.NoError(t, err)

	err = k.SetProof(ctx, types.Proof{
		MetadataUri: "ipfs://metadata2",
		Sender:      sender1.String(),
		Indices:     []int64{},
		Proofs:      [][]byte{},
		IsValidData: false,
	})
	require.NoError(t, err)

	err = k.SetProof(ctx, types.Proof{
		MetadataUri: "ipfs://metadata1",
		Sender:      sender2.String(),
		Indices:     []int64{},
		Proofs:      [][]byte{},
		IsValidData: false,
	})
	require.NoError(t, err)

	proof = k.GetProof(ctx, "ipfs://metadata1", sender1.String())
	require.Equal(t, proof.MetadataUri, "ipfs://metadata1")
	require.Equal(t, proof.Sender, sender1.String())
	require.Len(t, proof.Indices, 1)
	require.Len(t, proof.Proofs, 1)
	require.True(t, proof.IsValidData)

	proof = k.GetProof(ctx, "ipfs://metadata2", sender1.String())
	require.Equal(t, proof.MetadataUri, "ipfs://metadata2")
	require.Equal(t, proof.Sender, sender1.String())
	require.Len(t, proof.Indices, 0)
	require.Len(t, proof.Proofs, 0)
	require.False(t, proof.IsValidData)

	proof = k.GetProof(ctx, "ipfs://metadata1", sender2.String())
	require.Equal(t, proof.MetadataUri, "ipfs://metadata1")
	require.Equal(t, proof.Sender, sender2.String())
	require.Len(t, proof.Indices, 0)
	require.Len(t, proof.Proofs, 0)
	require.False(t, proof.IsValidData)

	proofs := k.GetProofs(ctx, "ipfs://metadata1")
	require.Len(t, proofs, 2)

	proofs = k.GetProofs(ctx, "ipfs://metadata2")
	require.Len(t, proofs, 1)

	proofs = k.GetAllProofs(ctx)
	require.Len(t, proofs, 3)

	k.DeleteProof(ctx, "ipfs://metadata1", sender1.String())
	proof = k.GetProof(ctx, "ipfs://metadata1", sender1.String())
	require.Equal(t, proof.MetadataUri, "")

	proofs = k.GetProofs(ctx, "ipfs://metadata1")
	require.Len(t, proofs, 1)

	proofs = k.GetAllProofs(ctx)
	require.Len(t, proofs, 2)
}
