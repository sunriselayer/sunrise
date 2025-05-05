package keeper

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k msgServer) BundleCommitments(ctx context.Context, msg *types.MsgBundleCommitments) (*types.MsgBundleCommitmentsResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	if len(msg.Commitments) != len(msg.Signatures) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "commitments and signatures must have the same length")
	}

	if len(msg.Commitments) == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "commitments must not be empty")
	}

	declarationHeight := msg.Commitments[0].DeclarationHeight
	shardsMerkleRoot := msg.Commitments[0].ShardsMerkleRoot

	declaration, found, err := k.GetBlobDeclaration(ctx, shardsMerkleRoot)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get blob declaration")
	}
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "blob declaration not found")
	}
	if declaration.BlockHeight != declarationHeight {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "declaration height mismatch")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if sdkCtx.BlockTime().After(declaration.Expiry) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "declaration has expired")
	}

	for i, commitment := range msg.Commitments {
		if commitment.DeclarationHeight != declarationHeight {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "all commitments must have the same declaration height")
		}

		if !bytes.Equal(commitment.ShardsMerkleRoot, shardsMerkleRoot) {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "all commitments must have the same shards merkle root")
		}

		validator, err := k.StakingKeeper.ValidatorAddressCodec().StringToBytes(commitment.Validator)
		if err != nil {
			return nil, errorsmod.Wrapf(err, "invalid validator address at index %d", i)
		}

		err = k.VerifyCommitmentSignature(ctx, validator, commitment, msg.Signatures[i])
		if err != nil {
			return nil, errorsmod.Wrapf(err, "invalid signature at index %d", i)
		}
	}

	err = k.TallyCommitments(ctx)
	if err != nil {
		return nil, err
	}

	if len(declaration.BundlerRewards) > 0 {
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, declaration.BundlerRewards)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to send bundler reward")
		}
	}

	return &types.MsgBundleCommitmentsResponse{}, nil
}
