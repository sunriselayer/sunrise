package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetProofDeputy(ctx context.Context, validator []byte) (deputy []byte, found bool) {
	has, err := k.ProofDeputies.Has(ctx, validator)
	if err != nil {
		panic(err)
	}

	if !has {
		return deputy, false
	}

	deputy, err = k.ProofDeputies.Get(ctx, validator)
	if err != nil {
		panic(err)
	}

	return deputy, true
}

// SetProofDeputy set the proof deputy of the validator
func (k Keeper) SetProofDeputy(ctx context.Context, validator []byte, deputy []byte) error {
	err := k.ProofDeputies.Set(ctx, validator, deputy)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteProofDeputy(ctx sdk.Context, validator []byte) {
	err := k.ProofDeputies.Remove(ctx, validator)
	if err != nil {
		panic(err)
	}
}
