package keeper

import (
	"bytes"
	"context"

	"cosmossdk.io/errors"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/cron/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// AddSchedule adds new schedule
func (k msgServer) AddSchedule(goCtx context.Context, req *types.MsgAddSchedule) (*types.MsgAddScheduleResponse, error) {
	authority, err := k.addressCodec.StringToBytes(req.Authority)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if !bytes.Equal(k.GetAuthority(), authority) {
		expectedAuthorityStr, _ := k.addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", expectedAuthorityStr, req.Authority)
	}

	if req.Name == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "name is invalid")
	}

	if req.Period == 0 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "period is invalid")
	}

	if len(req.Msgs) == 0 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "msgs should not be empty")
	}

	if _, ok := types.ExecutionStage_name[int32(req.ExecutionStage)]; !ok {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "execution stage is invalid")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.Keeper.AddSchedule(ctx, req.Name, req.Period, req.Msgs, req.ExecutionStage); err != nil {
		return nil, errors.Wrap(err, "failed to add schedule")
	}

	return &types.MsgAddScheduleResponse{}, nil
}

// RemoveSchedule removes schedule
func (k msgServer) RemoveSchedule(goCtx context.Context, req *types.MsgRemoveSchedule) (*types.MsgRemoveScheduleResponse, error) {
	authority, err := k.addressCodec.StringToBytes(req.Authority)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if !bytes.Equal(k.GetAuthority(), authority) {
		expectedAuthorityStr, _ := k.addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", expectedAuthorityStr, req.Authority)
	}

	if req.Name == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "name is invalid")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err = k.Keeper.RemoveSchedule(ctx, req.Name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove schedule")
	}

	return &types.MsgRemoveScheduleResponse{}, nil
}
