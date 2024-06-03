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

	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"

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

	maxAmountIn, ok := sdkmath.NewIntFromString(data.Amount)
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
	// TODO: validate converted denomIn is equal to the route DenomIn
	_ = denomIn

	var (
		result            types.RouteResult
		interfaceFee      sdkmath.Int
		remainderAmountIn sdkmath.Int
	)

	if metadata.ExactAmountIn != nil {
		// Swap exact amount in
		amountIn := maxAmountIn
		minAmountOut := metadata.ExactAmountIn.MinAmountOut

		result, interfaceFee, err = im.keeper.SwapExactAmountIn(
			ctx,
			swapper,
			metadata.InterfaceProvider,
			metadata.Route,
			amountIn,
			minAmountOut,
		)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}
	} else {
		// Swap exact amount out
		amountOut := metadata.ExactAmountOut.AmountOut

		result, interfaceFee, err = im.keeper.SwapExactAmountOut(
			ctx,
			swapper,
			metadata.InterfaceProvider,
			metadata.Route,
			maxAmountIn,
			amountOut,
		)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}
		amountIn := result.TokenIn.Amount

		remainderAmountIn = maxAmountIn.Sub(amountIn)

		if remainderAmountIn.IsPositive() && metadata.ExactAmountOut.Return == nil {
			// Send from swapper to receiver
			remainderTokenIn := sdk.NewCoin(denomIn, remainderAmountIn)
			err = im.keeper.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(remainderTokenIn))
			if err != nil {
				return channeltypes.NewErrorAcknowledgement(err)
			}

			remainderAmountIn = sdkmath.ZeroInt()
		}
	}

	// TODO: go to keeper and do this after forward packet finished
	// Return the remainder token in
	// err := im.keeper.TransferSwappedToken(
	// 	ctx,
	// 	swapper,
	// 	remainderTokenIn,
	// 	incomingAck.Acknowledgement(),
	// 	*metadata.ReturnAmountIn,
	// )
	// if err != nil {
	// 	return channeltypes.NewErrorAcknowledgement(err)
	// }

	denomOut := metadata.Route.DenomOut
	amountOut := result.TokenOut.Amount.Sub(interfaceFee)
	tokenOut := sdk.NewCoin(denomOut, amountOut)

	// Send from swapper to receiver
	err = im.keeper.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(tokenOut))
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	if metadata.Forward != nil {
		// Forward the swapped token out

		var returnMetadata *packetforwardtypes.ForwardMetadata = nil
		if metadata.ExactAmountOut != nil {
			returnMetadata = metadata.ExactAmountOut.Return
		}

		err := im.keeper.TransferSwappedToken(
			ctx,
			receiver.String(),
			tokenOut,
			*metadata.Forward,
			incomingAck.Acknowledgement(),
			result,
			remainderAmountIn,
			returnMetadata,
		)
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
		Result:      inflightPacket.Result,
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
			Result:      inflightPacket.Result,
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
