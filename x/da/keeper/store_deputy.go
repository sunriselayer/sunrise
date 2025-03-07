package keeper

import (
	"context"
)

func (k Keeper) GetProofDeputy(ctx context.Context, validator []byte) (deputy []byte, found bool, err error) {
	has, err := k.ProofDeputies.Has(ctx, validator)
	if err != nil {
		return deputy, false, err
	}

	if !has {
		return deputy, false, nil
	}

	deputy, err = k.ProofDeputies.Get(ctx, validator)
	if err != nil {
		return deputy, false, err
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

func (k Keeper) DeleteProofDeputy(ctx context.Context, validator []byte) error {
	err := k.ProofDeputies.Remove(ctx, validator)
	if err != nil {
		return err
	}

	return nil
}
