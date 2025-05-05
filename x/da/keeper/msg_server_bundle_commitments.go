package keeper

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
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

	declaration, found, err := k.GetBlobDeclaration(ctx, declarationHeight, shardsMerkleRoot)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get blob declaration")
	}
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "blob declaration not found")
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

		var pubKey cryptotypes.PubKey
		commitmentKey, found, err := k.GetCommitmentKey(ctx, validator)
		if err != nil {
			return nil, errorsmod.Wrapf(err, "failed to get commitment key for validator %s", validator)
		}
		if found {
		}

		signMessage, err := commitment.Marshal()
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to marshal commitment")
		}

		if !pubKey.VerifySignature(signMessage, msg.Signatures[i]) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid signature at index %d", i)
		}
	}

	return &types.MsgBundleCommitmentsResponse{}, nil
}
