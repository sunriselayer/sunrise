package swap

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v9/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v9/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v9/modules/core/05-port/types"
	exported "github.com/cosmos/ibc-go/v9/modules/core/exported"

	keeper "github.com/sunriselayer/sunrise/x/swap/keeper"
	types "github.com/sunriselayer/sunrise/x/swap/types"
)

var _ porttypes.IBCModule = IBCMiddleware{}

type IBCMiddleware struct {
	porttypes.IBCModule
	keeper *keeper.Keeper
}

// NewIBCMiddleware creates a new IBCMiddleware given the keeper and underlying application.
func NewIBCMiddleware(
	app porttypes.IBCModule,
	k *keeper.Keeper,
) IBCMiddleware {
	return IBCMiddleware{
		IBCModule: app,
		keeper:    k,
	}
}

// receiveFunds receives funds from the packet into the override receiver
// address and returns an error if the funds cannot be received.
func (im IBCMiddleware) receiveFunds(
	ctx context.Context,
	channelVersion string,
	packet channeltypes.Packet,
	data transfertypes.FungibleTokenPacketData,
	overrideReceiver string,
	relayer sdk.AccAddress,
) (exported.Acknowledgement, error) {
	overrideData := transfertypes.FungibleTokenPacketData{
		Denom:    data.Denom,
		Amount:   data.Amount,
		Sender:   data.Sender,
		Receiver: overrideReceiver, // override receiver
		// Memo explicitly zeroed
	}
	overrideDataBz := transfertypes.ModuleCdc.MustMarshalJSON(&overrideData)
	overridePacket := channeltypes.Packet{
		Sequence:           packet.Sequence,
		SourcePort:         packet.SourcePort,
		SourceChannel:      packet.SourceChannel,
		DestinationPort:    packet.DestinationPort,
		DestinationChannel: packet.DestinationChannel,
		Data:               overrideDataBz, // override data
		TimeoutHeight:      packet.TimeoutHeight,
		TimeoutTimestamp:   packet.TimeoutTimestamp,
	}

	ack := im.IBCModule.OnRecvPacket(ctx, channelVersion, overridePacket, relayer)

	if ack == nil {
		return ack, fmt.Errorf("ack is nil")
	}

	if !ack.Success() {
		return ack, fmt.Errorf("ack error: %s", string(ack.Acknowledgement()))
	}

	return ack, nil
}

// OnRecvPacket checks the memo field on this packet and if the metadata inside's root key indicates this packet
// should be handled by the swap middleware it attempts to perform a swap. If the swap is successful
// the underlying application's OnRecvPacket callback is invoked, an ack error is returned otherwise.
func (im IBCMiddleware) OnRecvPacket(
	ctx context.Context,
	channelVersion string,
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
		return im.IBCModule.OnRecvPacket(ctx, channelVersion, packet, relayer)
	}

	m, err := types.DecodeSwapMetadata(data.Memo)
	if err != nil {
		return im.IBCModule.OnRecvPacket(ctx, channelVersion, packet, relayer)
	}
	metadata := *m.Swap

	if err := metadata.Validate(); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	// Swap
	denomIn := types.GetDenomForThisChain(
		packet.DestinationPort,
		packet.DestinationChannel,
		packet.SourcePort,
		packet.SourceChannel,
		data.Denom,
	)
	if metadata.Route.DenomIn != denomIn {
		return channeltypes.NewErrorAcknowledgement(fmt.Errorf("invalid route: expected %s, got %s", metadata.Route.DenomIn, denomIn))
	}

	// Settle the incoming fund
	swapper := im.keeper.AccountKeeper.GetModuleAddress(types.ModuleName)
	incomingAck, err := im.receiveFunds(ctx, channelVersion, packet, data, swapper.String(), relayer)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	result, interfaceFee, err := im.keeper.SwapIncomingFund(
		sdkCtx,
		packet,
		swapper,
		data,
		metadata,
	)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	waitingPacket, err := im.keeper.ProcessSwappedFund(
		sdkCtx,
		packet,
		swapper,
		data,
		metadata,
		result,
		interfaceFee,
		incomingAck,
	)

	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	if waitingPacket != nil {
		err = im.keeper.SetIncomingInFlightPacket(ctx, *waitingPacket)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}

		// Returning nil ack will prevent WriteAcknowledgement from occurring for forwarded packet.
		// This is intentional so that the acknowledgement will be written later based on the ack/timeout of the forwarded packet.
		return nil
	}

	fullAck := types.SwapAcknowledgement{
		Result:      result,
		IncomingAck: incomingAck.Acknowledgement(),
	}

	bz, err := fullAck.Acknowledgement()
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	return channeltypes.NewResultAcknowledgement(bz)
}

// OnAcknowledgementPacket implements the IBCModule interface.
func (im IBCMiddleware) OnAcknowledgementPacket(
	ctx context.Context,
	channelVersion string,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return im.IBCModule.OnAcknowledgementPacket(ctx, channelVersion, packet, acknowledgement, relayer)
	}

	inflightPacket, found, err := im.keeper.GetOutgoingInFlightPacket(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence)
	if err != nil {
		return err
	}
	if !found {
		return im.IBCModule.OnAcknowledgementPacket(ctx, channelVersion, packet, acknowledgement, relayer)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = im.keeper.OnAcknowledgementOutgoingInFlightPacket(sdkCtx, packet, acknowledgement, inflightPacket)
	if err != nil {
		return err
	}

	return im.IBCModule.OnAcknowledgementPacket(ctx, channelVersion, packet, acknowledgement, relayer)
}

// OnTimeoutPacket implements the IBCModule interface.
func (im IBCMiddleware) OnTimeoutPacket(
	ctx context.Context,
	channelVersion string,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return im.IBCModule.OnTimeoutPacket(ctx, channelVersion, packet, relayer)
	}

	inflightPacket, found, err := im.keeper.GetOutgoingInFlightPacket(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence)
	if err != nil {
		return err
	}
	if !found {
		return im.IBCModule.OnTimeoutPacket(ctx, channelVersion, packet, relayer)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := im.keeper.OnTimeoutOutgoingInFlightPacket(sdkCtx, packet, inflightPacket); err != nil {
		return err
	}

	return im.IBCModule.OnTimeoutPacket(ctx, channelVersion, packet, relayer)
}
