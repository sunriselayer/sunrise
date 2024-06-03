package keeper

import (
	"context"
	"encoding/json"

	errors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/swap/types"

	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	exported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

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

	if swapData.ExactAmountIn != nil {
		// Swap exact amount in
		amountIn := maxAmountIn
		minAmountOut := swapData.ExactAmountIn.MinAmountOut

		result, interfaceFee, err = k.SwapExactAmountIn(
			ctx,
			swapper,
			swapData.InterfaceProvider,
			swapData.Route,
			amountIn,
			minAmountOut,
		)
		if err != nil {
			return types.RouteResult{}, sdkmath.Int{}, err
		}
	} else {
		// Swap exact amount out
		amountOut := swapData.ExactAmountOut.AmountOut

		result, interfaceFee, err = k.SwapExactAmountOut(
			ctx,
			swapper,
			swapData.InterfaceProvider,
			swapData.Route,
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
		Index:        types.NewPacketIndex(incomingPacket.DestinationPort, incomingPacket.DestinationChannel, incomingPacket.Sequence),
		Data:         incomingPacket.Data,
		SrcPortId:    incomingPacket.SourcePort,
		SrcChannelId: incomingPacket.SourceChannel,
		Ack:          incomingAck.Acknowledgement(),
		Result:       result,
		InterfaceFee: interfaceFee,
		Return:       &types.IncomingInFlightPacket_AckReturn{},  // default value is nil ack
		Forward:      &types.IncomingInFlightPacket_AckForward{}, // default value is nil ack
	}

	maxAmountIn, ok := sdkmath.NewIntFromString(tokenData.Amount)
	if !ok {
		return nil, errors.Wrap(sdkerrors.ErrInvalidCoins, "invalid amount")
	}
	remainderAmountIn := maxAmountIn.Sub(result.TokenIn.Amount)

	waiting := false

	if remainderAmountIn.IsPositive() {
		remainderTokenIn := sdk.NewCoin(result.TokenIn.Denom, remainderAmountIn)

		if swapData.ExactAmountOut.Return != nil {
			// Return the remainder token in
			returnPacket, err := k.TransferAndCreateOutgoingInFlightPacket(
				ctx,
				waitingPacket.Index,
				tokenData.Receiver,
				remainderTokenIn,
				*swapData.ExactAmountOut.Return,
			)
			if err != nil {
				return nil, err
			}

			waitingPacket.Return = &types.IncomingInFlightPacket_OutgoingIndexReturn{
				OutgoingIndexReturn: &returnPacket.Index,
			}
			waiting = true
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
	ctx context.Context,
	ackWaitingIndex types.PacketIndex,
	sender string,
	tokenOut sdk.Coin,
	metadata packetforwardtypes.ForwardMetadata,
) (packet types.OutgoingInFlightPacket, err error) {
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

	packet = types.OutgoingInFlightPacket{
		Index: types.NewPacketIndex(
			metadata.Port,
			metadata.Channel,
			res.Sequence,
		),
		AckWaitingIndex:  ackWaitingIndex,
		RetriesRemaining: int32(retries),
	}

	k.SetOutgoingInFlightPacket(ctx, packet)

	return packet, nil
}

func (k Keeper) OnAcknowledgementOutgoingInFlightPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	inFlightPacket types.OutgoingInFlightPacket,
) error {
	waitingPacket, found := k.GetIncomingInFlightPacket(ctx, inFlightPacket.AckWaitingIndex.PortId, inFlightPacket.AckWaitingIndex.ChannelId, inFlightPacket.AckWaitingIndex.Sequence)
	if found {
		k.RemoveOutgoingInFlightPacket(ctx, inFlightPacket.Index.PortId, inFlightPacket.Index.ChannelId, inFlightPacket.Index.Sequence)
	} else {
		return nil
	}

	// The pattern of waitingPacket.Return == nil is not handled here
	switch waitingPacket.Return.(type) {
	case *types.IncomingInFlightPacket_OutgoingIndexReturn:
		waitingPacket.Return = &types.IncomingInFlightPacket_AckReturn{
			AckReturn: acknowledgement,
		}
		break
	case *types.IncomingInFlightPacket_AckReturn:
		break
	}

	// The pattern of waitingPacket.Forward == nil is not handled here
	switch waitingPacket.Forward.(type) {
	case *types.IncomingInFlightPacket_OutgoingIndexForward:
		waitingPacket.Forward = &types.IncomingInFlightPacket_AckForward{
			AckForward: acknowledgement,
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

	return nil
}

func (k Keeper) OnTimeoutOutgoingInFlightPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	inFlightPacket types.OutgoingInFlightPacket,
) error {
	k.RemoveOutgoingInFlightPacket(ctx, inFlightPacket.Index.PortId, inFlightPacket.Index.ChannelId, inFlightPacket.Index.Sequence)
	inFlightPacket.RetriesRemaining--

	if inFlightPacket.RetriesRemaining > 0 {
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
			nil, // TODO
			nil, // TODO
			packet.Data,
		)
		if err != nil {
			return err
		}

		// Set the new sequence number
		inFlightPacket.Index.Sequence = sequence
		k.SetOutgoingInFlightPacket(ctx, inFlightPacket)
	} else {
		// If remaining retry count is zero:
		// - Returning non error acknowledgement to the origin
		// - However it contains error acknowledgement of forward / return packet
		ack := channeltypes.NewErrorAcknowledgement(errors.Wrap(sdkerrors.ErrUnknownRequest, "Retry count on timeout exceeds"))

		waitingPacket, found := k.GetIncomingInFlightPacket(ctx, inFlightPacket.AckWaitingIndex.PortId, inFlightPacket.AckWaitingIndex.ChannelId, inFlightPacket.AckWaitingIndex.Sequence)
		if !found {
			return nil
		}

		switch packetReturn := waitingPacket.Return.(type) {
		case *types.IncomingInFlightPacket_OutgoingIndexReturn:
			if packetReturn.OutgoingIndexReturn.Equal(inFlightPacket.Index) {
				waitingPacket.Return = &types.IncomingInFlightPacket_AckReturn{
					AckReturn: ack.Acknowledgement(),
				}
			}
			break
		case *types.IncomingInFlightPacket_AckReturn:
			break
		}

		switch packetForward := waitingPacket.Forward.(type) {
		case *types.IncomingInFlightPacket_OutgoingIndexForward:
			if packetForward.OutgoingIndexForward.Equal(inFlightPacket.Index) {
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
	switch packet.Return.(type) {
	case *types.IncomingInFlightPacket_OutgoingIndexReturn:
		return false, nil
	case *types.IncomingInFlightPacket_AckReturn:
		break
	}

	switch packet.Forward.(type) {
	case *types.IncomingInFlightPacket_OutgoingIndexForward:
		return false, nil
	case *types.IncomingInFlightPacket_AckForward:
		break
	}

	fullAck := types.SwapAcknowledgement{
		Result:      packet.Result,
		IncomingAck: packet.Ack,
		ReturnAck:   packet.Return.(*types.IncomingInFlightPacket_AckReturn).AckReturn,
		ForwardAck:  packet.Forward.(*types.IncomingInFlightPacket_AckForward).AckForward,
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
			nil, // TODO
			nil, // TODO
		),
		channeltypes.NewResultAcknowledgement(bz),
	); err != nil {
		return false, err
	}

	k.RemoveIncomingInFlightPacket(ctx, packet.Index.PortId, packet.Index.ChannelId, packet.Index.Sequence)

	return true, nil
}
