package keeper

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"

	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func (k Keeper) TransferAndCreateInFlightPacket(
	ctx context.Context,
	sender string,
	tokenOut sdk.Coin,
	metadata packetforwardtypes.ForwardMetadata,
) (packet types.InFlightPacket, err error) {
	var memo string
	if metadata.Next != nil {
		if err := json.Unmarshal([]byte(memo), &metadata.Next); err != nil {
			return packet, err
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
		return packet, err
	}

	var retries uint8
	if metadata.Retries != nil {
		retries = *metadata.Retries
	} else {
		retries = types.DefaultRetryCount
	}

	packet = types.InFlightPacket{
		Index: types.NewPacketIndex(
			metadata.Port,
			metadata.Channel,
			res.Sequence,
		),
		RetriesRemaining: int32(retries),
	}

	k.SetInFlightPacket(ctx, packet)

	return packet, nil
}

func (k Keeper) OnAcknowledgementInFlightPacket() {}

func (k Keeper) OnTimeoutInFlightPacket() {}
