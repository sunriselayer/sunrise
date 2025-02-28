package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetProofDeputy(ctx context.Context, validator []byte) (deputy []byte, found bool, err error) {
	has, err := k.ProofDeputies.Has(ctx, validator)
	if err != nil {
		return nil, false, err
	}

	if !has {
		return nil, false, nil
	}

	deputy, err = k.ProofDeputies.Get(ctx, validator)
	if err != nil {
		return nil, false, err
	}

	return deputy, true, nil
}

// SetProofDeputy set the proof deputy of the validator
func (k Keeper) SetProofDeputy(ctx context.Context, validator []byte, deputy []byte) error {
	err := k.ProofDeputies.Set(ctx, validator, deputy)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteProofDeputy(ctx sdk.Context, validator []byte) error {
	return k.ProofDeputies.Remove(ctx, validator)
}
