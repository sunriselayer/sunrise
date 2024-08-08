package keeper

import (
	"context"
)

func (k Keeper) EndBlocker(ctx context.Context) {
	// TODO: after proof period, remove proofs and slash validators not submitted proofs
	// TODO: send slashed funds to challenger
}
