package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) RegisterCommitmentKey(ctx context.Context, msg *types.MsgRegisterCommitmentKey) (*types.MsgRegisterCommitmentKeyResponse, error) {
	validator, err := k.StakingKeeper.ValidatorAddressCodec().StringToBytes(msg.Validator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	err = k.SetCommitmentKey(ctx, validator, types.CommitmentKey{
		Validator: msg.Validator,
		Pubkey:    msg.Pubkey,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterCommitmentKeyResponse{}, nil
}

func (k msgServer) UnregisterCommitmentKey(ctx context.Context, msg *types.MsgUnregisterCommitmentKey) (*types.MsgUnregisterCommitmentKeyResponse, error) {
	validator, err := k.StakingKeeper.ValidatorAddressCodec().StringToBytes(msg.Validator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	err = k.DeleteCommitmentKey(ctx, validator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to delete commitment key")
	}

	return &types.MsgUnregisterCommitmentKeyResponse{}, nil
}
