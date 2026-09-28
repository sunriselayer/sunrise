package keeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrIfShutdownInFlightLeavesOrdinaryPacketsAlone(t *testing.T) {
	inactive := Keeper{}
	err := inactive.ErrIfShutdownInFlight(context.Background(), true, "OnTimeoutPacket", "transfer", "channel-0", 7)
	require.NoError(t, err)

	active := Keeper{ShutdownActive: func(context.Context) bool { return true }}
	err = active.ErrIfShutdownInFlight(context.Background(), false, "OnTimeoutPacket", "transfer", "channel-0", 7)
	require.NoError(t, err)

	err = active.ErrIfShutdownInFlight(context.Background(), true, "OnTimeoutPacket", "transfer", "channel-0", 7)
	require.Error(t, err)
	require.Contains(t, err.Error(), "OnTimeoutPacket")
	require.Contains(t, err.Error(), "transfer/channel-0/7")
}

func TestIsShutdownActive(t *testing.T) {
	require.False(t, Keeper{}.IsShutdownActive(context.Background()))
	keeper := Keeper{ShutdownActive: func(context.Context) bool { return true }}
	require.True(t, keeper.IsShutdownActive(context.Background()))
}
