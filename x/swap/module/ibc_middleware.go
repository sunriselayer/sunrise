package swap

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	exported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	keeper "github.com/sunriselayer/sunrise/x/swap/keeper"
	types "github.com/sunriselayer/sunrise/x/swap/types"
)

type IBCMiddleware struct {
	porttypes.IBCModule
	keeper keeper.Keeper
}

// NewIBCMiddleware creates a new IBCMiddleware given the keeper and underlying application.
func NewIBCMiddleware(
	app porttypes.IBCModule,
	k keeper.Keeper,
) IBCMiddleware {
	return IBCMiddleware{
		IBCModule: app,
		keeper:    k,
	}
}

// receiveFunds receives funds from the packet into the override receiver
// address and returns an error if the funds cannot be received.
func (im IBCMiddleware) receiveFunds(
	ctx sdk.Context,
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

	ack := im.IBCModule.OnRecvPacket(ctx, overridePacket, relayer)

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
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	d := make(map[string]interface{})
	err := json.Unmarshal([]byte(data.Memo), &d)
	if err != nil || d["swap"] == nil {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	m := &types.PacketMetadata{}
	err = json.Unmarshal([]byte(data.Memo), m)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(fmt.Errorf("error parsing swap metadata: %w", err))
	}

	metadata := m.Swap

	if err := metadata.Validate(); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	// Prepare for swap
	swapper := im.keeper.AccountKeeper.GetModuleAddress(types.ModuleName)
	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	amountIn, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		return channeltypes.NewErrorAcknowledgement(errors.Wrap(sdkerrors.ErrInvalidCoins, "invalid amount"))
	}

	// Settle the incoming fund
	incomingAck, err := im.receiveFunds(ctx, packet, data, swapper.String(), relayer)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	// Swap
	denomIn := data.Denom // TODO: convert ibc denom
	// TODO: validaatep converted denomIn is equal to the route DenomIn
	_ = denomIn

	amountOut, err := im.keeper.SwapExactAmountIn(
		ctx,
		swapper,
		metadata.InterfaceProvider,
		metadata.Route,
		amountIn,
		metadata.MinAmountOut,
	)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}
	denomOut := metadata.Route.DenomOut
	tokenOut := sdk.NewCoin(denomOut, amountOut)

	// Transfer from swapper to receiver
	err = im.keeper.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(tokenOut))
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	if metadata.Forward != nil {
		err := im.keeper.ForwardSwappedToken(
			ctx,
			swapper,
			tokenOut,
			incomingAck.Acknowledgement(),
			*metadata.Forward,
		)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}

		// returning nil ack will prevent WriteAcknowledgement from occurring for forwarded packet.
		// This is intentional so that the acknowledgement will be written later based on the ack/timeout of the forwarded packet.
		return nil
	}

	fullAck := types.SwapAcknowledgement{
		AmountOut:   amountOut,
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
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return im.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	inflightPacket, found := im.keeper.GetInFlightPacket(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence)
	if !found {
		return im.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	fullAck := types.SwapAcknowledgement{
		AmountOut:   inflightPacket.AmountOut,
		IncomingAck: inflightPacket.IncomingAck,
		ForwardAck:  acknowledgement,
	}
	bz, err := fullAck.Acknowledgement()
	if err != nil {
		return err
	}
	im.keeper.RemoveInFlightPacket(ctx, inflightPacket.SrcPortId, inflightPacket.SrcChannelId, inflightPacket.Sequence)

	return im.IBCModule.OnAcknowledgementPacket(ctx, packet, bz, relayer)
}

// OnTimeoutPacket implements the IBCModule interface.
func (im IBCMiddleware) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return im.IBCModule.OnTimeoutPacket(ctx, packet, relayer)
	}

	inflightPacket, found := im.keeper.GetInFlightPacket(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence)
	if !found {
		return im.IBCModule.OnTimeoutPacket(ctx, packet, relayer)
	}

	inflightPacket.RetriesRemaining--

	if inflightPacket.RetriesRemaining > 0 {
		// TODO: Resend packet

		im.keeper.RemoveInFlightPacket(ctx, inflightPacket.SrcPortId, inflightPacket.SrcChannelId, inflightPacket.Sequence)
		inflightPacket.Sequence = 0 // TODO: Reset sequence
		im.keeper.SetInFlightPacket(ctx, inflightPacket)
	} else {
		// If remaining retry count is zero:
		// - Returning non error acknowledgement to the origin
		// - However it contains error acknowledgement of forwarding packet
		forwardAck := channeltypes.NewErrorAcknowledgement(errors.Wrap(sdkerrors.ErrUnknownRequest, "Retry count on timeout exceeds"))
		fullAck := types.SwapAcknowledgement{
			AmountOut:   inflightPacket.AmountOut,
			IncomingAck: inflightPacket.IncomingAck,
			ForwardAck:  forwardAck.Acknowledgement(),
		}
		bz, err := fullAck.Acknowledgement()
		if err != nil {
			return err
		}

		if err := im.keeper.IbcKeeperFn().ChannelKeeper.WriteAcknowledgement(
			ctx,
			nil, // TODO
			nil, // TODO
			channeltypes.NewResultAcknowledgement(bz),
		); err != nil {
			return err
		}

		im.keeper.RemoveInFlightPacket(ctx, inflightPacket.SrcPortId, inflightPacket.SrcChannelId, inflightPacket.Sequence)
	}

	return im.IBCModule.OnTimeoutPacket(ctx, packet, relayer)
}
