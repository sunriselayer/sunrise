package keeper

import (
	"context"
	"fmt"
)

// IsShutdownActive reports whether incoming IBC should be rejected and in-flight
// packet completion should be left untouched.
func (k Keeper) IsShutdownActive(ctx context.Context) bool {
	if k.ShutdownActive == nil {
		return false
	}
	return k.ShutdownActive(ctx)
}

// ErrIfShutdownInFlight returns an error when a recorded in-flight packet is
// delivered after the shutdown upgrade. The caller must not refund, resend, or
// call the IBC keeper. Packets that are not in the in-flight index are unaffected.
func (k Keeper) ErrIfShutdownInFlight(ctx context.Context, found bool, caller, portID, channelID string, sequence uint64) error {
	if !found || !k.IsShutdownActive(ctx) {
		return nil
	}
	return fmt.Errorf("%s: in-flight packet %s/%s/%d is left in escrow and is not completed", caller, portID, channelID, sequence)
}
