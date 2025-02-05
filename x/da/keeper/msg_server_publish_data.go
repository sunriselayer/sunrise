package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) PublishData(ctx context.Context, msg *types.MsgPublishData) (*types.MsgPublishDataResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	if msg.ParityShardCount >= uint64(len(msg.ShardDoubleHashes)) {
		return nil, types.ErrParityShardCountGTETotal
	}
	// end static validation

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	err = k.SetPublishedData(ctx, types.PublishedData{
		Publisher:         msg.Sender,
		MetadataUri:       msg.MetadataUri,
		ParityShardCount:  msg.ParityShardCount,
		ShardDoubleHashes: msg.ShardDoubleHashes,
		Collateral:        params.PublishDataCollateral,
		Timestamp:         sdkCtx.BlockTime(),
		DataSourceInfo:    msg.DataSourceInfo,
		Status:            types.Status_STATUS_VOTING,
	})
	if err != nil {
		return nil, err
	}

	// Send collateral to module account
	if params.PublishDataCollateral.IsAllPositive() {
		sender := sdk.MustAccAddressFromBech32(msg.Sender)
		err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, params.PublishDataCollateral)
		if err != nil {
			return nil, err
		}
	}

	err = sdkCtx.EventManager().EmitTypedEvent(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgPublishDataResponse{}, nil
}
