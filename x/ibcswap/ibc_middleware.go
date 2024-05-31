package ibcswap

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	exported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

const ModuleName = "ibcswap"

type ibcSwapMiddleware struct {
	porttypes.IBCModule
}

func NewIBCMiddleware(ibcModule porttypes.IBCModule) porttypes.IBCModule {
	return &ibcSwapMiddleware{
		IBCModule: ibcModule,
	}
}

func (m *ibcSwapMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		// If this happens either a) a user has crafted an invalid packet, b) a
		// software developer has connected the middleware to a stack that does
		// not have a transfer module, or c) the transfer module has been modified
		// to accept other Packets. The best thing we can do here is pass the packet
		// on down the stack.
		return m.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	swapSuccess := false
	if swapSuccess {
		return m.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	ackErr := errors.Wrapf(sdkerrors.ErrInvalidType, "Failed to swap")

	return channeltypes.NewErrorAcknowledgement(ackErr)
}
