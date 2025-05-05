package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	types "github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) VerifyCommitmentSignature(ctx context.Context, validator sdk.ValAddress, commitment types.AvailabilityCommitment, signature []byte) error {
	var pubkey cryptotypes.PubKey

	commitmentKey, found, err := k.GetCommitmentKey(ctx, validator)
	if err != nil {
		return errorsmod.Wrapf(err, "failed to get commitment key for validator %s", validator)
	}
	if !found {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "commitment key not found for validator %s", validator)
	}

	// TODO: unmarshal pubkey

	signMessage, err := commitment.Marshal()
	if err != nil {
		return errorsmod.Wrap(err, "failed to marshal commitment")
	}

	if !pubkey.VerifySignature(signMessage, signature) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid signature")
	}

	return nil
}

func (k Keeper) TallyCommitments(ctx context.Context) error {

}

func (k Keeper) TakeValidatorsPowerSnapshot(ctx context.Context) error {
	iterator, err := k.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		return err
	}

	defer iterator.Close()

	snapshot := types.ValidatorsPowerSnapshot{
		BlockHeight: sdk.UnwrapSDKContext(ctx).BlockHeight(),
		Snapshots:   []types.ValidatorPowerSnapshot{},
	}

	for ; iterator.Valid(); iterator.Next() {
		validator, err := k.StakingKeeper.Validator(ctx, iterator.Value())
		if err != nil {
			return err
		}

		power := validator.GetBondedTokens()
		powerSnapshot := types.ValidatorPowerSnapshot{
			Validator: validator.GetOperator(),
			Power:     power,
		}

		snapshot.Snapshots = append(snapshot.Snapshots, powerSnapshot)
	}

	return k.ValidatorsPowerSnapshots.Set(ctx, snapshot.BlockHeight, snapshot)
}
