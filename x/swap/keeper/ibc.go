package keeper

import (
	"time"

	errors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/swap/types"

	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	exported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var (
	// DefaultTransferPacketTimeoutHeight is the timeout height following IBC defaults
	DefaultTransferPacketTimeoutHeight = clienttypes.Height{
		RevisionNumber: 0,
		RevisionHeight: 0,
	}

	// DefaultTransferPacketTimeoutTimestamp is the timeout timestamp following IBC defaults
	DefaultTransferPacketTimeoutTimestamp = time.Duration(transfertypes.DefaultRelativePacketTimeoutTimestamp) * time.Nanosecond
)

func timeoutTimestamp(ctx sdk.Context, duration time.Duration) uint64 {
	return uint64(ctx.BlockTime().UnixNano()) + uint64(duration.Nanoseconds())
}

func (k Keeper) SwapIncomingFund(
	ctx sdk.Context,
	incomingPacket channeltypes.Packet,
	swapper sdk.AccAddress,
	tokenData transfertypes.FungibleTokenPacketData,
	swapData types.SwapMetadata,
) (result types.RouteResult, interfaceFee sdkmath.Int, err error) {
	maxAmountIn, ok := sdkmath.NewIntFromString(tokenData.Amount)
	if !ok {
		return result, interfaceFee, errors.Wrap(sdkerrors.ErrInvalidCoins, "invalid amount")
	}

	// Prepare for swap
	receiver, err := sdk.AccAddressFromBech32(tokenData.Receiver)
	if err != nil {
		return result, interfaceFee, err
	}

	switch amountStrategy := swapData.AmountStrategy.(type) {
	case *types.SwapMetadata_ExactAmountIn:
		// Swap exact amount in
		amountIn := maxAmountIn
		minAmountOut := amountStrategy.ExactAmountIn.MinAmountOut

		result, interfaceFee, err = k.SwapExactAmountIn(
			ctx,
			swapper,
			swapData.InterfaceProvider,
			*swapData.Route,
			amountIn,
			minAmountOut,
		)
		if err != nil {
			return types.RouteResult{}, sdkmath.Int{}, err
		}
	case *types.SwapMetadata_ExactAmountOut:
		// Swap exact amount out
		amountOut := amountStrategy.ExactAmountOut.AmountOut

		result, interfaceFee, err = k.SwapExactAmountOut(
			ctx,
			swapper,
			swapData.InterfaceProvider,
			*swapData.Route,
			maxAmountIn,
			amountOut,
		)
		if err != nil {
			return types.RouteResult{}, sdkmath.Int{}, err
		}
	}

	denomOut := swapData.Route.DenomOut
	amountOutGross := result.TokenOut.Amount
	amountOutNet := amountOutGross.Sub(interfaceFee)

	tokenOutNet := sdk.NewCoin(denomOut, amountOutNet)

	// Send from swapper to receiver
	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(tokenOutNet))
	if err != nil {
		return types.RouteResult{}, sdkmath.Int{}, err
	}

	return result, interfaceFee, nil
}

func (k Keeper) ProcessSwappedFund(
	ctx sdk.Context,
	incomingPacket channeltypes.Packet,
	swapper sdk.AccAddress,
	tokenData transfertypes.FungibleTokenPacketData,
	swapData types.SwapMetadata,
	result types.RouteResult,
	interfaceFee sdkmath.Int,
	incomingAck exported.Acknowledgement,
) (waitingPacket *types.IncomingInFlightPacket, err error) {
	waitingPacket = &types.IncomingInFlightPacket{
		Index:            types.NewPacketIndex(incomingPacket.DestinationPort, incomingPacket.DestinationChannel, incomingPacket.Sequence),
		Data:             incomingPacket.Data,
		SrcPortId:        incomingPacket.SourcePort,
		SrcChannelId:     incomingPacket.SourceChannel,
		TimeoutHeight:    incomingPacket.TimeoutHeight.String(),
		TimeoutTimestamp: incomingPacket.TimeoutTimestamp,
		Ack:              incomingAck.Acknowledgement(),
		Result:           result,
		InterfaceFee:     interfaceFee,
		Change:           &types.IncomingInFlightPacket_AckChange{},  // default value is nil ack
		Forward:          &types.IncomingInFlightPacket_AckForward{}, // default value is nil ack
	}

	maxAmountIn, ok := sdkmath.NewIntFromString(tokenData.Amount)
	if !ok {
		return nil, errors.Wrap(sdkerrors.ErrInvalidCoins, "invalid amount")
	}
	remainderAmountIn := maxAmountIn.Sub(result.TokenIn.Amount)

	waiting := false

	if remainderAmountIn.IsPositive() {
		remainderTokenIn := sdk.NewCoin(result.TokenIn.Denom, remainderAmountIn)

		switch amountStrategy := swapData.AmountStrategy.(type) {
		case *types.SwapMetadata_ExactAmountOut:
			if amountStrategy.ExactAmountOut.Change != nil {
				// Return the remainder token in
				returnPacket, err := k.TransferAndCreateOutgoingInFlightPacket(
					ctx,
					waitingPacket.Index,
					tokenData.Receiver,
					remainderTokenIn,
					*amountStrategy.ExactAmountOut.Change,
				)
				if err != nil {
					return nil, err
				}

				waitingPacket.Change = &types.IncomingInFlightPacket_OutgoingIndexChange{
					OutgoingIndexChange: &returnPacket.Index,
				}
				waiting = true
			}
		}
	}

	if swapData.Forward != nil {
		amountOutGross := result.TokenOut.Amount
		amountOutNet := amountOutGross.Sub(interfaceFee)

		tokenOutNet := sdk.NewCoin(result.TokenOut.Denom, amountOutNet)

		// Forward the swapped token out
		forwardPacket, err := k.TransferAndCreateOutgoingInFlightPacket(
			ctx,
			waitingPacket.Index,
			tokenData.Receiver,
			tokenOutNet,
			*swapData.Forward,
		)
		if err != nil {
			return nil, err
		}

		waitingPacket.Forward = &types.IncomingInFlightPacket_OutgoingIndexForward{
			OutgoingIndexForward: &forwardPacket.Index,
		}
		waiting = true
	}

	if waiting {
		k.SetIncomingInFlightPacket(ctx, *waitingPacket)

		return waitingPacket, nil
	}

	return nil, nil
}

func (k Keeper) TransferAndCreateOutgoingInFlightPacket(
	ctx sdk.Context,
	incomingIndex types.PacketIndex,
	sender string,
	tokenOut sdk.Coin,
	metadata types.ForwardMetadata,
) (packet types.OutgoingInFlightPacket, err error) {

	msgTransfer := transfertypes.MsgTransfer{
		SourcePort:       metadata.Port,
		SourceChannel:    metadata.Channel,
		Token:            tokenOut,
		Sender:           sender,
		Receiver:         metadata.Receiver,
		TimeoutHeight:    DefaultTransferPacketTimeoutHeight,
		TimeoutTimestamp: timeoutTimestamp(ctx, metadata.Timeout),
		Memo:             metadata.Next,
	}
	// forward token to receiver
	res, err := k.TransferKeeper.Transfer(ctx, &msgTransfer)
	if err != nil {
		return packet, err
	}

	var retries uint8
	if metadata.Retries == 0 {
		retries = types.DefaultRetryCount
	} else {
		retries = uint8(metadata.Retries)
	}

	packet = types.OutgoingInFlightPacket{
		Index: types.NewPacketIndex(
			metadata.Port,
			metadata.Channel,
			res.Sequence,
		),
		AckWaitingIndex:  incomingIndex,
		RetriesRemaining: int32(retries),
	}

	k.SetOutgoingInFlightPacket(ctx, packet)

	return packet, nil
}

func (k Keeper) OnAcknowledgementOutgoingInFlightPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	outgoingPacket types.OutgoingInFlightPacket,
) error {
	incomingPacket, found := k.GetIncomingInFlightPacket(ctx, outgoingPacket.AckWaitingIndex.PortId, outgoingPacket.AckWaitingIndex.ChannelId, outgoingPacket.AckWaitingIndex.Sequence)
	if found {
		k.RemoveOutgoingInFlightPacket(ctx, outgoingPacket.Index.PortId, outgoingPacket.Index.ChannelId, outgoingPacket.Index.Sequence)
	} else {
		return nil
	}

	// The pattern of waitingPacket.Return == nil is not handled here
	switch incomingPacket.Change.(type) {
	case *types.IncomingInFlightPacket_OutgoingIndexChange:
		incomingPacket.Change = &types.IncomingInFlightPacket_AckChange{
			AckChange: acknowledgement,
		}
		break
	case *types.IncomingInFlightPacket_AckChange:
		break
	}

	// The pattern of waitingPacket.Forward == nil is not handled here
	switch incomingPacket.Forward.(type) {
	case *types.IncomingInFlightPacket_OutgoingIndexForward:
		incomingPacket.Forward = &types.IncomingInFlightPacket_AckForward{
			AckForward: acknowledgement,
		}
		break
	case *types.IncomingInFlightPacket_AckForward:
		break
	}

	deleted, err := k.ShouldDeleteCompletedWaitingPacket(ctx, incomingPacket)
	if err != nil {
		return err
	}
	if !deleted {
		k.SetIncomingInFlightPacket(ctx, incomingPacket)
	}

	return nil
}

func (k Keeper) OnTimeoutOutgoingInFlightPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	outgoingPacket types.OutgoingInFlightPacket,
) error {
	k.RemoveOutgoingInFlightPacket(ctx, outgoingPacket.Index.PortId, outgoingPacket.Index.ChannelId, outgoingPacket.Index.Sequence)
	outgoingPacket.RetriesRemaining--

	if outgoingPacket.RetriesRemaining > 0 {
		// Resend packet
		_, chanCap, err := k.IbcKeeperFn().ChannelKeeper.LookupModuleByChannel(ctx, packet.DestinationPort, packet.DestinationChannel)
		if err != nil {
			return errors.Wrap(err, "could not retrieve module from port-id")
		}
		sequence, err := k.IbcKeeperFn().ChannelKeeper.SendPacket(
			ctx,
			chanCap,
			packet.SourcePort,
			packet.SourceChannel,
			DefaultTransferPacketTimeoutHeight,
			timeoutTimestamp(ctx, DefaultTransferPacketTimeoutTimestamp),
			packet.Data,
		)
		if err != nil {
			return err
		}

		// Set the new sequence number
		outgoingPacket.Index.Sequence = sequence
		k.SetOutgoingInFlightPacket(ctx, outgoingPacket)
	} else {
		// If remaining retry count is zero:
		// - Returning non error acknowledgement to the origin
		// - However it contains error acknowledgement of change / forward packet
		ack := channeltypes.NewErrorAcknowledgement(errors.Wrap(sdkerrors.ErrUnknownRequest, "Retry count on timeout exceeds"))

		waitingPacket, found := k.GetIncomingInFlightPacket(ctx, outgoingPacket.AckWaitingIndex.PortId, outgoingPacket.AckWaitingIndex.ChannelId, outgoingPacket.AckWaitingIndex.Sequence)
		if !found {
			return nil
		}

		switch packetReturn := waitingPacket.Change.(type) {
		case *types.IncomingInFlightPacket_OutgoingIndexChange:
			if packetReturn.OutgoingIndexChange.Equal(outgoingPacket.Index) {
				waitingPacket.Change = &types.IncomingInFlightPacket_AckChange{
					AckChange: ack.Acknowledgement(),
				}
			}
			break
		case *types.IncomingInFlightPacket_AckChange:
			break
		}

		switch packetForward := waitingPacket.Forward.(type) {
		case *types.IncomingInFlightPacket_OutgoingIndexForward:
			if packetForward.OutgoingIndexForward.Equal(outgoingPacket.Index) {
				waitingPacket.Forward = &types.IncomingInFlightPacket_AckForward{
					AckForward: ack.Acknowledgement(),
				}
			}
			break
		case *types.IncomingInFlightPacket_AckForward:
			break
		}

		deleted, err := k.ShouldDeleteCompletedWaitingPacket(ctx, waitingPacket)
		if err != nil {
			return err
		}
		if !deleted {
			k.SetIncomingInFlightPacket(ctx, waitingPacket)
		}
	}

	return nil
}

func (k Keeper) ShouldDeleteCompletedWaitingPacket(
	ctx sdk.Context,
	packet types.IncomingInFlightPacket,
) (deleted bool, err error) {
	switch packet.Change.(type) {
	case *types.IncomingInFlightPacket_OutgoingIndexChange:
		return false, nil
	case *types.IncomingInFlightPacket_AckChange:
		break
	}

	switch packet.Forward.(type) {
	case *types.IncomingInFlightPacket_OutgoingIndexForward:
		return false, nil
	case *types.IncomingInFlightPacket_AckForward:
		break
	}

	var changeAck []byte = nil
	var forwardAck []byte = nil

	if packet.Change != nil {
		changeAck = packet.Change.(*types.IncomingInFlightPacket_AckChange).AckChange
	}

	if packet.Forward != nil {
		forwardAck = packet.Forward.(*types.IncomingInFlightPacket_AckForward).AckForward
	}

	fullAck := types.SwapAcknowledgement{
		Result:      packet.Result,
		IncomingAck: packet.Ack,
		ChangeAck:   changeAck,
		ForwardAck:  forwardAck,
	}
	bz, err := fullAck.Acknowledgement()
	if err != nil {
		return false, err
	}

	_, chanCap, err := k.IbcKeeperFn().ChannelKeeper.LookupModuleByChannel(ctx, packet.Index.PortId, packet.Index.ChannelId)
	if err != nil {
		return false, errors.Wrap(err, "could not retrieve module from port-id")
	}
	if err := k.IbcKeeperFn().ChannelKeeper.WriteAcknowledgement(
		ctx,
		chanCap,
		channeltypes.NewPacket(
			packet.Data,
			packet.Index.Sequence,
			packet.SrcPortId,
			packet.SrcChannelId,
			packet.Index.PortId,
			packet.Index.ChannelId,
			clienttypes.MustParseHeight(packet.TimeoutHeight),
			packet.TimeoutTimestamp,
		),
		channeltypes.NewResultAcknowledgement(bz),
	); err != nil {
		return false, err
	}

	k.RemoveIncomingInFlightPacket(ctx, packet.Index.PortId, packet.Index.ChannelId, packet.Index.Sequence)

	return true, nil
}
