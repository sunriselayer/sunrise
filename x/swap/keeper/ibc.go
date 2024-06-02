package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"

	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func (k Keeper) ForwardSwappedToken(
	ctx context.Context,
	swapper sdk.AccAddress,
	token sdk.Coin,
	metadata packetforwardtypes.ForwardMetadata,
) error {
	msgTransfer := transfertypes.MsgTransfer{
		SourcePort: metadata.Port,
	}
	// forward token to receiver
	res, err := k.transferKeeper.Transfer(ctx, &msgTransfer)
	if err != nil {
		return err
	}

	k.SetInFlightPacket(ctx, types.InFlightPacket{
		SrcPortId:    metadata.Port,
		SrcChannelId: metadata.Channel,
		Sequence:     res.Sequence,
	})

	return nil
}
