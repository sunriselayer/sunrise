package keeper_test

import (
	"bufio"
	"bytes"
	"math/big"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/consensys/gnark-crypto/ecc"
	native_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/sunriselayer/sunrise/x/da/types"
	"github.com/sunriselayer/sunrise/x/da/zkp"
)

func TestMsgSubmitProof(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	params := types.DefaultParams()
	require.NoError(t, k.SetParams(ctx, params))
	wctx := sdk.UnwrapSDKContext(ctx)

	preImage1 := big.NewInt(111)
	m := native_mimc.NewMiMC()
	m.Write(preImage1.Bytes())
	hash := m.Sum(nil)

	// witness definition
	assignment := zkp.ValidityProofCircuit{
		ShardHash:       preImage1,
		ShardDoubleHash: hash,
	}

	// compiles our circuit into a R1CS
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &zkp.ValidityProofCircuit{})
	require.NoError(t, err)

	// Recover proving key
	provingKey, err := zkp.UnmarshalProvingKey(params.ZkpProvingKey)
	require.NoError(t, err)

	witness1, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	require.NoError(t, err)

	// groth16: Prove & Verify
	proof, err := groth16.Prove(ccs, provingKey, witness1)
	require.NoError(t, err)

	var b bytes.Buffer
	bufWrite := bufio.NewWriter(&b)
	_, err = proof.WriteTo(bufWrite)
	require.NoError(t, err)
	err = bufWrite.Flush()
	require.NoError(t, err)
	proofBytes := b.Bytes()

	err = k.SetPublishedData(ctx, types.PublishedData{
		MetadataUri:        "ipfs://metadata1",
		ParityShardCount:   0,
		ShardDoubleHashes:  [][]byte{hash},
		Timestamp:          time.Time{},
		Status:             "challenge_for_fraud",
		Publisher:          "publisher",
		Challenger:         "challenger",
		Collateral:         sdk.Coins{},
		ChallengeTimestamp: time.Time{},
	})
	require.NoError(t, err)

	testCases := []struct {
		name      string
		input     *types.MsgSubmitProof
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid data hash",
			input: &types.MsgSubmitProof{
				Sender:      "sender",
				MetadataUri: "ipfs://metadata1",
				Indices:     []int64{},
				Proofs:      [][]byte{},
				IsValidData: false,
			},
			expErr: false,
		},
		{
			name: "valid proof",
			input: &types.MsgSubmitProof{
				Sender:      "sender",
				MetadataUri: "ipfs://metadata1",
				Indices:     []int64{0},
				Proofs:      [][]byte{proofBytes},
				IsValidData: true,
			},
			expErr: false,
		},
		{
			name: "invalid proof",
			input: &types.MsgSubmitProof{
				Sender:      "sender",
				MetadataUri: "ipfs://metadata1",
				Indices:     []int64{0},
				Proofs:      [][]byte{{0x0}},
				IsValidData: true,
			},
			expErr:    true,
			expErrMsg: "unexpected EOF",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.SubmitProof(wctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
