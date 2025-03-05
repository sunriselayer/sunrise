package keeper_test

import (
	"bufio"
	"bytes"
	"math/big"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/consensys/gnark-crypto/ecc"
	native_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/sunriselayer/sunrise/x/da/types"
	"github.com/sunriselayer/sunrise/x/da/zkp"
)

type mockValidator struct {
	sdk.ValidatorI
}

func (mv mockValidator) IsBonded() bool {
	return true
}

func TestMsgSubmitValidityProof(t *testing.T) {
	k, mocks, ms, ctx := setupMsgServer(t)
	params := types.DefaultParams()
	require.NoError(t, k.Params.Set(ctx, params))
	sender := sdk.AccAddress("sender")
	validator := sdk.ValAddress("validator")
	validatorAcc := sdk.AccAddress(validator)
	mockVal := mockValidator{}
	mocks.StakingKeeper.EXPECT().Validator(gomock.Any(), gomock.Any()).Return(mockVal, nil).AnyTimes()
	_, err := ms.RegisterProofDeputy(ctx, &types.MsgRegisterProofDeputy{
		Sender:        validatorAcc.String(),
		DeputyAddress: sender.String(),
	})
	require.NoError(t, err)
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
		MetadataUri:                "ipfs://metadata1",
		ParityShardCount:           0,
		ShardDoubleHashes:          [][]byte{hash},
		Timestamp:                  time.Now(),
		Status:                     types.Status_STATUS_CHALLENGING,
		Publisher:                  "publisher",
		Challenger:                 "challenger",
		PublishDataCollateral:      sdk.Coins{},
		SubmitInvalidityCollateral: sdk.Coins{},
		ChallengeTimestamp:         time.Now(),
	})
	require.NoError(t, err)

	testCases := []struct {
		name      string
		input     *types.MsgSubmitValidityProof
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid validator address",
			input: &types.MsgSubmitValidityProof{
				Sender:           sender.String(),
				ValidatorAddress: sdk.ValAddress("invalid").String(),
				MetadataUri:      "ipfs://metadata1",
				Indices:          []int64{},
				Proofs:           [][]byte{},
			},
			expErr:    true,
			expErrMsg: "proof deputy not found",
		},
		{
			name: "empty proof",
			input: &types.MsgSubmitValidityProof{
				Sender:           sender.String(),
				ValidatorAddress: validator.String(),
				MetadataUri:      "ipfs://metadata1",
				Indices:          []int64{},
				Proofs:           [][]byte{},
			},
			expErr: false,
		},
		{
			name: "valid proof",
			input: &types.MsgSubmitValidityProof{
				Sender:           sender.String(),
				ValidatorAddress: validator.String(),
				MetadataUri:      "ipfs://metadata1",
				Indices:          []int64{0},
				Proofs:           [][]byte{proofBytes},
			},
			expErr: false,
		},
		{
			name: "invalid proof",
			input: &types.MsgSubmitValidityProof{
				Sender:           sender.String(),
				ValidatorAddress: validator.String(),
				MetadataUri:      "ipfs://metadata1",
				Indices:          []int64{0},
				Proofs:           [][]byte{{0x0}},
			},
			expErr:    true,
			expErrMsg: "unexpected EOF",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.SubmitValidityProof(ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
