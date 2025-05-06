package keeper

import (
	"context"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/da/das/consts"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k msgServer) DeclareBlob(ctx context.Context, msg *types.MsgDeclareBlob) (*types.MsgDeclareBlobResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	if msg.MetadataSize == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "metadata size must be positive")
	}
	if len(msg.ShardsMerkleRoot) != 32 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "shards merkle root must be 32 bytes poseidon hash")
	}
	if msg.Rows == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "rows must be positive")
	}
	if msg.Cols == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cols must be positive")
	}
	// cols must be 2^n
	if (msg.Cols & (msg.Cols - 1)) != 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cols must be 2^n")
	}
	// rows > 1, cols must be 2*srs length
	if msg.Rows > 1 && msg.Cols != consts.ExtensionRatio*consts.SrsLen {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "if rows > 1, cols must be %d", consts.ExtensionRatio*consts.SrsLen)
	}
	if len(msg.KzgCommitmentsMerkleRoot) != 32 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "kzg commitments merkle root must be 32 bytes poseidon hash")
	}
	if err := msg.BundlerRewards.Validate(); err != nil {
		return nil, errorsmod.Wrap(err, "invalid bundler reward, empty is possible")
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Consume gas
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.GasMeter().ConsumeGas(params.GasPerByte*msg.MetadataSize, "declare blob metadata size")
	sdkCtx.GasMeter().ConsumeGas(params.GasPerByte*types.CalculateShardsTotalSize(msg.ShardCount), "declare blob size")

	_, has, err := k.GetBlobDeclaration(ctx, msg.ShardsMerkleRoot)
	if err != nil {
		return nil, err
	}
	if has {
		return nil, types.ErrDeclarationAlreadyExists
	}

	err = k.SetBlobDeclaration(ctx, types.BlobDeclaration{
		Sender:                   msg.Sender,
		BlockHeight:              sdkCtx.BlockHeight(),
		ShardsMerkleRoot:         msg.ShardsMerkleRoot,
		ShardCount:               msg.ShardCount,
		KzgCommitmentsMerkleRoot: msg.KzgCommitmentsMerkleRoot,
		BundlerRewards:           msg.BundlerRewards,
		Expiry:                   sdkCtx.BlockTime().Add(params.DeclarationPeriod),
	})
	if err != nil {
		return nil, err
	}

	// Send collateral to module account
	if len(msg.BundlerRewards) > 0 {
		err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, msg.BundlerRewards)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgDeclareBlobResponse{}, nil
}
