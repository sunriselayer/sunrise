package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
)

func TestForwardTimeoutDurationUsesDefaultWhenUnset(t *testing.T) {
	require.Equal(t, DefaultTransferPacketTimeoutTimestamp, forwardTimeoutDuration(0))
	require.Equal(t, DefaultTransferPacketTimeoutTimestamp, forwardTimeoutDuration(-time.Second))

	custom := 2 * time.Minute
	require.Equal(t, custom, forwardTimeoutDuration(custom))
	require.Equal(t, 10*time.Minute, DefaultTransferPacketTimeoutTimestamp)
}

func TestGetIBCKeeperRejectsNilCallback(t *testing.T) {
	swapKeeper := Keeper{}

	ibcKeeper, err := swapKeeper.getIBCKeeper()
	require.Nil(t, ibcKeeper)
	require.Error(t, err)
	require.ErrorContains(t, err, "IbcKeeperFn is nil")

	swapKeeper.IbcKeeperFn = func() *ibckeeper.Keeper { return nil }
	ibcKeeper, err = swapKeeper.getIBCKeeper()
	require.Nil(t, ibcKeeper)
	require.Error(t, err)
	require.ErrorContains(t, err, "returned a nil IBC keeper")
}
