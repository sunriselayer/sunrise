package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func (k Keeper) ForwardSwappedToken(
	ctx context.Context,
	swapper sdk.AccAddress,
	token sdk.Coin,
	metadata packetforwardtypes.ForwardMetadata,
) error {
	msgTransfer := transfertypes.MsgTransfer{}
	// forward token to receiver
	res, err := k.transferKeeper.Transfer(ctx, &msgTransfer)
	if err != nil {
		return err
	}

	// TODO: save forwarding packet
	_ = res.Sequence

	return nil
}
