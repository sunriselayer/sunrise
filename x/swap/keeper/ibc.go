package keeper

import (
	stderrors "errors"
	"time"

	errors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	exported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

// forwardTimeoutDuration returns the relative timeout for a swap forward.
// A zero duration is treated as unset, the same way retries of zero use
// DefaultRetryCount. Adding a zero duration to the block time creates a
// packet that is already expired, so it can never be received.
func forwardTimeoutDuration(timeout time.Duration) time.Duration {
	if timeout <= 0 {
		return DefaultTransferPacketTimeoutTimestamp
	}
	return timeout
}

// getIBCKeeper returns the app IBC keeper. Depinject builds the swap keeper
// before the IBC keeper exists, so the app must assign IbcKeeperFn afterwards.
// Calling a nil callback panics inside acknowledgement and timeout handling.
func (k Keeper) getIBCKeeper() (*ibckeeper.Keeper, error) {
	if k.IbcKeeperFn == nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "swap IbcKeeperFn is nil; the app must set it after the IBC keeper is created")
	}
	ibcKeeper := k.IbcKeeperFn()
	if ibcKeeper == nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "swap IbcKeeperFn returned a nil IBC keeper")
	}
	return ibcKeeper, nil
}

var (
	// DefaultTransferPacketTimeoutHeight is the timeout height following IBC defaults
	DefaultTransferPacketTimeoutHeight = clienttypes.Height{
		RevisionNumber: 0,
		RevisionHeight: 0,
	}

	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
	// DefaultTransferPacketTimeoutTimestamp is the timeout timestamp following IBC defaults
	DefaultTransferPacketTimeoutTimestamp = time.Duration(DefaultRelativePacketTimeoutTimestamp) * time.Nanosecond
)

func timeoutTimestamp(ctx sdk.Context, duration time.Duration) uint64 {
	return uint64(ctx.BlockTime().UnixNano()) + uint64(duration.Nanoseconds())
}

// resendTimedOutSwapPacket sends the refunded tokens again after a swap
// forward times out. The caller must already have run the transfer module's
// OnTimeoutPacket, which returns the escrowed coins to the original sender.
func (k Keeper) resendTimedOutSwapPacket(ctx sdk.Context, packet channeltypes.Packet) (uint64, error) {
	if k.TransferKeeper == nil {
		return 0, errors.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"transfer keeper is nil, cannot resend timed-out swap packet %s/%s/%d",
			packet.SourcePort, packet.SourceChannel, packet.Sequence,
		)
	}

	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return 0, errors.Wrapf(
			err,
			"failed to unmarshal fungible token packet data for %s/%s/%d",
			packet.SourcePort, packet.SourceChannel, packet.Sequence,
		)
	}

	amount, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		return 0, errors.Wrapf(
			sdkerrors.ErrInvalidCoins,
			"invalid amount %q on timed-out swap packet %s/%s/%d",
			data.Amount, packet.SourcePort, packet.SourceChannel, packet.Sequence,
		)
	}
	// The packet denom is the ICS-20 path (e.g. transfer/channel-0/uusdc), not
	// a bank denom. The transfer module refunds the timed-out packet as
	// ExtractDenomFromPath(path).IBCDenom(), which is ibc/{hash} for vouchers
	// and the base denom for native tokens, so resend that same denom.
	denom := transfertypes.ExtractDenomFromPath(data.Denom).IBCDenom()
	coin := sdk.Coin{Denom: denom, Amount: amount}
	if err := coin.Validate(); err != nil {
		return 0, errors.Wrapf(
			err,
			"invalid token %s%s (packet denom %s) on timed-out swap packet %s/%s/%d",
			data.Amount, denom, data.Denom, packet.SourcePort, packet.SourceChannel, packet.Sequence,
		)
	}

	res, err := k.TransferKeeper.Transfer(ctx, &transfertypes.MsgTransfer{
		SourcePort:       packet.SourcePort,
		SourceChannel:    packet.SourceChannel,
		Token:            coin,
		Sender:           data.Sender,
		Receiver:         data.Receiver,
		TimeoutHeight:    DefaultTransferPacketTimeoutHeight,
		TimeoutTimestamp: timeoutTimestamp(ctx, DefaultTransferPacketTimeoutTimestamp),
		Memo:             data.Memo,
	})
	if err != nil {
		return 0, errors.Wrapf(
			err,
			"failed to resend timed-out swap packet %s/%s/%d from sender %s to receiver %s",
			packet.SourcePort, packet.SourceChannel, packet.Sequence, data.Sender, data.Receiver,
		)
	}
	return res.Sequence, nil
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
		err = k.SetIncomingInFlightPacket(ctx, *waitingPacket)
		if err != nil {
			return nil, err
		}

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
		TimeoutTimestamp: timeoutTimestamp(ctx, forwardTimeoutDuration(metadata.Timeout)),
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

	err = k.SetOutgoingInFlightPacket(ctx, packet)
	if err != nil {
		return packet, err
	}

	return packet, nil
}

func (k Keeper) OnAcknowledgementOutgoingInFlightPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	outgoingPacket types.OutgoingInFlightPacket,
) error {
	incomingPacket, found, err := k.GetIncomingInFlightPacket(ctx, outgoingPacket.AckWaitingIndex.PortId, outgoingPacket.AckWaitingIndex.ChannelId, outgoingPacket.AckWaitingIndex.Sequence)
	if err != nil {
		return err
	}
	if found {
		err = k.RemoveOutgoingInFlightPacket(ctx, outgoingPacket.Index.PortId, outgoingPacket.Index.ChannelId, outgoingPacket.Index.Sequence)
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	// The pattern of waitingPacket.Return == nil is not handled here
	switch t := incomingPacket.Change.(type) {
	case *types.IncomingInFlightPacket_OutgoingIndexChange:
		if t != nil && t.OutgoingIndexChange != nil && t.OutgoingIndexChange.Equal(outgoingPacket.Index) {
			incomingPacket.Change = &types.IncomingInFlightPacket_AckChange{
				AckChange: acknowledgement,
			}
		}
		// case *types.IncomingInFlightPacket_AckChange:
	}

	// The pattern of waitingPacket.Forward == nil is not handled here
	switch t := incomingPacket.Forward.(type) {
	case *types.IncomingInFlightPacket_OutgoingIndexForward:
		if t != nil && t.OutgoingIndexForward != nil && t.OutgoingIndexForward.Equal(outgoingPacket.Index) {
			incomingPacket.Forward = &types.IncomingInFlightPacket_AckForward{
				AckForward: acknowledgement,
			}
		}
		// case *types.IncomingInFlightPacket_AckForward:
	}

	deleted, err := k.ShouldDeleteCompletedWaitingPacket(ctx, incomingPacket)
	if err != nil {
		return err
	}
	if !deleted {
		err = k.SetIncomingInFlightPacket(ctx, incomingPacket)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) OnTimeoutOutgoingInFlightPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	outgoingPacket types.OutgoingInFlightPacket,
) error {
	err := k.RemoveOutgoingInFlightPacket(ctx, outgoingPacket.Index.PortId, outgoingPacket.Index.ChannelId, outgoingPacket.Index.Sequence)
	if err != nil {
		return err
	}
	outgoingPacket.RetriesRemaining--

	if outgoingPacket.RetriesRemaining > 0 {
		// The transfer module refunds the timed-out escrow before this function
		// runs. Resend through MsgTransfer so the refunded coins are escrowed
		// again. ChannelKeeper.SendPacket only writes a commitment and would
		// leave that new packet unbacked. It also panics when IbcKeeperFn was
		// never wired, which is what blocked sequence 9631 on sunrise-1.
		sequence, err := k.resendTimedOutSwapPacket(ctx, packet)
		if err != nil {
			return err
		}

		outgoingPacket.Index.Sequence = sequence
		err = k.SetOutgoingInFlightPacket(ctx, outgoingPacket)
		if err != nil {
			return err
		}
	} else {
		// If remaining retry count is zero:
		// - Returning non error acknowledgement to the origin
		// - However it contains error acknowledgement of change / forward packet
		ack := channeltypes.NewErrorAcknowledgement(errors.Wrap(sdkerrors.ErrUnknownRequest, "Retry count on timeout exceeds"))

		waitingPacket, found, err := k.GetIncomingInFlightPacket(ctx, outgoingPacket.AckWaitingIndex.PortId, outgoingPacket.AckWaitingIndex.ChannelId, outgoingPacket.AckWaitingIndex.Sequence)
		if err != nil {
			return err
		}
		if !found {
			return nil
		}

		switch packetReturn := waitingPacket.Change.(type) {
		case *types.IncomingInFlightPacket_OutgoingIndexChange:
			if packetReturn != nil && packetReturn.OutgoingIndexChange != nil && packetReturn.OutgoingIndexChange.Equal(outgoingPacket.Index) {
				waitingPacket.Change = &types.IncomingInFlightPacket_AckChange{
					AckChange: ack.Acknowledgement(),
				}
			}
			// case *types.IncomingInFlightPacket_AckChange:
		}

		switch packetForward := waitingPacket.Forward.(type) {
		case *types.IncomingInFlightPacket_OutgoingIndexForward:
			if packetForward != nil && packetForward.OutgoingIndexForward != nil && packetForward.OutgoingIndexForward.Equal(outgoingPacket.Index) {
				waitingPacket.Forward = &types.IncomingInFlightPacket_AckForward{
					AckForward: ack.Acknowledgement(),
				}
			}
			// case *types.IncomingInFlightPacket_AckForward:
		}

		deleted, err := k.ShouldDeleteCompletedWaitingPacket(ctx, waitingPacket)
		if err != nil {
			return err
		}
		if !deleted {
			err = k.SetIncomingInFlightPacket(ctx, waitingPacket)
			if err != nil {
				return err
			}
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
	case *types.IncomingInFlightPacket_AckChange, nil:
		break
	}

	switch packet.Forward.(type) {
	case *types.IncomingInFlightPacket_OutgoingIndexForward:
		return false, nil
	case *types.IncomingInFlightPacket_AckForward, nil:
		break
	}

	var changeAck []byte = nil
	var forwardAck []byte = nil

	if ack, ok := packet.Change.(*types.IncomingInFlightPacket_AckChange); ok && ack != nil {
		changeAck = ack.AckChange
	}

	if ack, ok := packet.Forward.(*types.IncomingInFlightPacket_AckForward); ok && ack != nil {
		forwardAck = ack.AckForward
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

	ibcKeeper, err := k.getIBCKeeper()
	if err != nil {
		return false, errors.Wrapf(err, "cannot write acknowledgement for incoming packet %s/%s/%d", packet.Index.PortId, packet.Index.ChannelId, packet.Index.Sequence)
	}
	if err := ibcKeeper.ChannelKeeper.WriteAcknowledgement(
		ctx,
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
	); err != nil && !stderrors.Is(err, channeltypes.ErrAcknowledgementExists) {
		return false, err
	}

	err = k.RemoveIncomingInFlightPacket(ctx, packet.Index.PortId, packet.Index.ChannelId, packet.Index.Sequence)
	if err != nil {
		return false, err
	}

	return true, nil
}
