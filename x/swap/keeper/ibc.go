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
) (waitingPacket *types.AckWaitingPacket, err error) {
	waitingPacket = &types.AckWaitingPacket{
		Index:   types.NewPacketIndex(incomingPacket.DestinationPort, incomingPacket.DestinationChannel, incomingPacket.Sequence),
		Result:  result,
		Ack:     incomingAck.Acknowledgement(),
		Return:  &types.AckWaitingPacket_AckReturn{},  // default value is nil ack
		Forward: &types.AckWaitingPacket_AckForward{}, // default value is nil ack
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
			returnPacket, err := k.TransferAndCreateInFlightPacket(
				ctx,
				waitingPacket.Index,
				tokenData.Receiver,
				remainderTokenIn,
				*swapData.ExactAmountOut.Return,
			)
			if err != nil {
				return nil, err
			}

			waitingPacket.Return = &types.AckWaitingPacket_InFlightIndexReturn{
				InFlightIndexReturn: &returnPacket.Index,
			}
			waiting = true
		}
	}

	if swapData.Forward != nil {
		amountOutGross := result.TokenOut.Amount
		amountOutNet := amountOutGross.Sub(interfaceFee)

		tokenOutNet := sdk.NewCoin(result.TokenOut.Denom, amountOutNet)

		// Forward the swapped token out
		forwardPacket, err := k.TransferAndCreateInFlightPacket(
			ctx,
			waitingPacket.Index,
			tokenData.Receiver,
			tokenOutNet,
			*swapData.Forward,
		)
		if err != nil {
			return nil, err
		}

		waitingPacket.Forward = &types.AckWaitingPacket_InFlightIndexForward{
			InFlightIndexForward: &forwardPacket.Index,
		}
		waiting = true
	}

	if waiting {
		k.SetAckWaitingPacket(ctx, *waitingPacket)

		return waitingPacket, nil
	}

	return nil, nil
}

func (k Keeper) TransferAndCreateInFlightPacket(
	ctx context.Context,
	ackWaitingIndex types.PacketIndex,
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
		AckWaitingIndex:  ackWaitingIndex,
		RetriesRemaining: int32(retries),
	}

	k.SetInFlightPacket(ctx, packet)

	return packet, nil
}

func (k Keeper) OnAcknowledgementInFlightPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	inFlightPacket types.InFlightPacket,
) error {
	waitingPacket, found := k.GetAckWaitingPacket(ctx, inFlightPacket.AckWaitingIndex.PortId, inFlightPacket.AckWaitingIndex.ChannelId, inFlightPacket.AckWaitingIndex.Sequence)
	if found {
		k.RemoveInFlightPacket(ctx, inFlightPacket.Index.PortId, inFlightPacket.Index.ChannelId, inFlightPacket.Index.Sequence)
	} else {
		return nil
	}

	// The pattern of waitingPacket.Return == nil is not handled here
	switch waitingPacket.Return.(type) {
	case *types.AckWaitingPacket_InFlightIndexReturn:
		waitingPacket.Return = &types.AckWaitingPacket_AckReturn{
			AckReturn: acknowledgement,
		}
		break
	case *types.AckWaitingPacket_AckReturn:
		break
	}

	// The pattern of waitingPacket.Forward == nil is not handled here
	switch waitingPacket.Forward.(type) {
	case *types.AckWaitingPacket_InFlightIndexForward:
		waitingPacket.Forward = &types.AckWaitingPacket_AckForward{
			AckForward: acknowledgement,
		}
		break
	case *types.AckWaitingPacket_AckForward:
		break
	}

	deleted, err := k.ShouldDeleteCompletedWaitingPacket(ctx, waitingPacket)
	if err != nil {
		return err
	}
	if !deleted {
		k.SetAckWaitingPacket(ctx, waitingPacket)
	}

	return nil
}

func (k Keeper) OnTimeoutInFlightPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	inFlightPacket types.InFlightPacket,
) error {
	k.RemoveInFlightPacket(ctx, inFlightPacket.Index.PortId, inFlightPacket.Index.ChannelId, inFlightPacket.Index.Sequence)
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
		k.SetInFlightPacket(ctx, inFlightPacket)
	} else {
		// If remaining retry count is zero:
		// - Returning non error acknowledgement to the origin
		// - However it contains error acknowledgement of forward / return packet
		ack := channeltypes.NewErrorAcknowledgement(errors.Wrap(sdkerrors.ErrUnknownRequest, "Retry count on timeout exceeds"))

		waitingPacket, found := k.GetAckWaitingPacket(ctx, inFlightPacket.AckWaitingIndex.PortId, inFlightPacket.AckWaitingIndex.ChannelId, inFlightPacket.AckWaitingIndex.Sequence)
		if !found {
			return nil
		}

		switch packetReturn := waitingPacket.Return.(type) {
		case *types.AckWaitingPacket_InFlightIndexReturn:
			if packetReturn.InFlightIndexReturn.Equal(inFlightPacket.Index) {
				waitingPacket.Return = &types.AckWaitingPacket_AckReturn{
					AckReturn: ack.Acknowledgement(),
				}
			}
			break
		case *types.AckWaitingPacket_AckReturn:
			break
		}

		switch packetForward := waitingPacket.Forward.(type) {
		case *types.AckWaitingPacket_InFlightIndexForward:
			if packetForward.InFlightIndexForward.Equal(inFlightPacket.Index) {
				waitingPacket.Forward = &types.AckWaitingPacket_AckForward{
					AckForward: ack.Acknowledgement(),
				}
			}
			break
		case *types.AckWaitingPacket_AckForward:
			break
		}

		deleted, err := k.ShouldDeleteCompletedWaitingPacket(ctx, waitingPacket)
		if err != nil {
			return err
		}
		if !deleted {
			k.SetAckWaitingPacket(ctx, waitingPacket)
		}
	}

	return nil
}

func (k Keeper) ShouldDeleteCompletedWaitingPacket(
	ctx sdk.Context,
	packet types.AckWaitingPacket,
) (deleted bool, err error) {
	switch packet.Return.(type) {
	case *types.AckWaitingPacket_InFlightIndexReturn:
		return false, nil
	case *types.AckWaitingPacket_AckReturn:
		break
	}

	switch packet.Forward.(type) {
	case *types.AckWaitingPacket_InFlightIndexForward:
		return false, nil
	case *types.AckWaitingPacket_AckForward:
		break
	}

	fullAck := types.SwapAcknowledgement{
		Result:      packet.Result,
		IncomingAck: packet.Ack,
		ReturnAck:   packet.Return.(*types.AckWaitingPacket_AckReturn).AckReturn,
		ForwardAck:  packet.Forward.(*types.AckWaitingPacket_AckForward).AckForward,
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
			nil, // TODO
			packet.Index.Sequence,
			nil, // TODO
			nil, // TODO
			nil, // TODO
			nil, // TODO
			nil, // TODO
			nil, // TODO
		),
		channeltypes.NewResultAcknowledgement(bz),
	); err != nil {
		return false, err
	}

	k.RemoveAckWaitingPacket(ctx, packet.Index.PortId, packet.Index.ChannelId, packet.Index.Sequence)

	return true, nil
}
