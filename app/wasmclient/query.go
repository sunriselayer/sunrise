package wasmclient

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CustomQuery struct {
	VerifyGroth16Bn254    *QueryVerifyGroth16 `json:"verify_groth16_bn254,omitempty"`
	VerifyGroth16Bls12381 *QueryVerifyGroth16 `json:"verify_groth16_bls12381,omitempty"`
}

type QueryVerifyGroth16 struct {
	Proof         []byte `json:"proof"`
	VerifyingKey  []byte `json:"verifying_key"`
	PublicWitness []byte `json:"public_witness"`
}

func CustomQuerier() func(sdk.Context, json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var customQuery CustomQuery
		err := json.Unmarshal([]byte(request), &customQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to parse custom query %v", err)
		}

		switch {
		case customQuery.VerifyGroth16Bn254 != nil:
			err := VerifyGroth16WithCurve(
				customQuery.VerifyGroth16Bn254.Proof,
				customQuery.VerifyGroth16Bn254.VerifyingKey,
				customQuery.VerifyGroth16Bn254.PublicWitness,
				CurveBN254,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to verify groth16 BN254 proof: %v", err)
			}
			return []byte{}, nil
		case customQuery.VerifyGroth16Bls12381 != nil:
			err := VerifyGroth16WithCurve(
				customQuery.VerifyGroth16Bls12381.Proof,
				customQuery.VerifyGroth16Bls12381.VerifyingKey,
				customQuery.VerifyGroth16Bls12381.PublicWitness,
				CurveBLS12381,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to verify groth16 BLS12-381 proof: %v", err)
			}
			return []byte{}, nil
		default:
			return nil, fmt.Errorf("unknown custom query %v", request)
		}
	}

}
