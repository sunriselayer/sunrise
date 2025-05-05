package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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
	if msg.BlobSize == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "blob size must be positive")
	}
	if len(msg.KzgCommitmentsMerkleRoot) != 32 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "kzg commitments merkle root must be 32 bytes poseidon hash")
	}
	if msg.BundlerReward != nil {
		if err := msg.BundlerReward.Validate(); err != nil {
			return nil, errorsmod.Wrap(err, "invalid bundler reward, empty is possible")
		}
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Consume gas
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.GasMeter().ConsumeGas(params.GasPerByte*msg.MetadataSize, "declare blob metadata size")
	sdkCtx.GasMeter().ConsumeGas(params.GasPerByte*msg.BlobSize, "declare blob size")

	_, has, err := k.GetBlobDeclaration(ctx, sdkCtx.BlockHeight(), msg.ShardsMerkleRoot)
	if err != nil {
		return nil, err
	}
	if has {
		return nil, types.ErrDataAlreadyExist
	}

	err = k.SetBlobDeclaration(ctx, types.BlobDeclaration{
		Sender:                   msg.Sender,
		BlockHeight:              sdkCtx.BlockHeight(),
		ShardsMerkleRoot:         msg.ShardsMerkleRoot,
		ShardCount:               types.CalculateShardCount(msg.BlobSize),
		KzgCommitmentsMerkleRoot: msg.KzgCommitmentsMerkleRoot,
		BundlerReward:            msg.BundlerReward,
		Expiry:                   sdkCtx.BlockTime().Add(params.DeclarationPeriod),
	})
	if err != nil {
		return nil, err
	}

	// Send collateral to module account
	if msg.BundlerReward != nil {
		err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(sender), types.ModuleName, sdk.NewCoins(*msg.BundlerReward))
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgDeclareBlobResponse{}, nil
}
