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

	// check duplicate data
	has, err := k.PublishedData.Has(ctx, msg.MetadataUri)
	if err != nil {
		return nil, err
	}
	if has {
		return nil, types.ErrDataAlreadyExist
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Consume gas
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.GasMeter().ConsumeGas(params.PublishDataGas, "publish data")

	err = k.SetPublishedData(ctx, types.PublishedData{
		MetadataUri:                msg.MetadataUri,
		ParityShardCount:           msg.ParityShardCount,
		ShardDoubleHashes:          msg.ShardDoubleHashes,
		Timestamp:                  sdkCtx.BlockTime(),
		Status:                     types.Status_STATUS_CHALLENGE_PERIOD,
		Publisher:                  msg.Sender,
		PublishDataCollateral:      params.PublishDataCollateral,
		SubmitInvalidityCollateral: params.SubmitInvalidityCollateral,
		PublishedTimestamp:         sdkCtx.BlockTime(),
		DataSourceInfo:             msg.DataSourceInfo,
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
