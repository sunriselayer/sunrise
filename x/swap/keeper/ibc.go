package keeper

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"

	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func (k Keeper) TransferSwappedToken(
	ctx context.Context,
	swapper sdk.AccAddress,
	token sdk.Coin,
	incomingAck []byte,
	metadata packetforwardtypes.ForwardMetadata,
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
		Token:         token,
		Sender:        swapper.String(),
		Receiver:      metadata.Receiver,
		// TODO: timeout
		Memo: memo,
	}
	// forward token to receiver
	res, err := k.transferKeeper.Transfer(ctx, &msgTransfer)
	if err != nil {
		return err
	}

	var retries uint8
	if metadata.Retries != nil {
		retries = *metadata.Retries
	} else {
		retries = types.DefaultRetryCount
	}

	k.SetInFlightPacket(ctx, types.InFlightPacket{
		SrcPortId:        metadata.Port,
		SrcChannelId:     metadata.Channel,
		Sequence:         res.Sequence,
		RetriesRemaining: int32(retries),
		AmountOut:        token.Amount,
		IncomingAck:      incomingAck,
	})

	return nil
}
