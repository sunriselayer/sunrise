package keeper

import (
	"context"
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"

	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func (k Keeper) TransferSwappedToken(
	ctx context.Context,
	sender string,
	tokenOut sdk.Coin,
	metadata packetforwardtypes.ForwardMetadata,
	incomingAck []byte,
	result types.RouteResult,
	remainderAmountIn sdkmath.Int,
	returnMetadata *packetforwardtypes.ForwardMetadata,
) error {
	var memo string
	if metadata.Next != nil {
		if err := json.Unmarshal([]byte(memo), &metadata.Next); err != nil {
			return err
		}
	}

	msgTransfer := transfertypes.MsgTransfer{
		SourcePort:    metadata.Port,
		SourceChannel: metadata.Channel,
		Token:         tokenOut,
		Sender:        sender,
		Receiver:      metadata.Receiver,
		// TODO: timeout
		Memo: memo,
	}
	// forward token to receiver
	res, err := k.TransferKeeper.Transfer(ctx, &msgTransfer)
	if err != nil {
		return err
	}

	var retries uint8
	if metadata.Retries != nil {
		retries = *metadata.Retries
	} else {
		retries = types.DefaultRetryCount
	}

	val := types.InFlightPacket{
		SrcPortId:        metadata.Port,
		SrcChannelId:     metadata.Channel,
		Sequence:         res.Sequence,
		RetriesRemaining: int32(retries),
		IncomingAck:      incomingAck,
		Result:           result,
	}

	if returnMetadata != nil {
		val.ReturnInfo = &types.ReturnInfo{
			// TODO
			RemainderAmountIn: remainderAmountIn,
		}
	}

	k.SetInFlightPacket(ctx, val)

	return nil
}
