package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

// The following code is a new implementation of the MsgServer interface
// for the da module. It handles the processing of MsgVerifyData messages.
// This was previously handled by the EndBlocker.
func (k msgServer) VerifyData(goCtx context.Context, msg *types.MsgVerifyData) (*types.MsgVerifyDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// If STATUS_REJECTED is overtime, remove from the store
	err = k.DeleteRejectedDataOvertime(ctx, params.RejectedRemovalPeriod)
	if err != nil {
		return nil, err
	}

	// IF STATUS_VERIFIED is overtime, remove from store
	err = k.DeleteVerifiedDataOvertime(ctx, params.VerifiedRemovalPeriod)
	if err != nil {
		return nil, err
	}

	// if STATUS_CHALLENGE_PERIOD receives invalidity above the threshold, change to STATUS_CHALLENGING
	err = k.ChangeToChallengingFromChallengePeriod(ctx, params.ChallengeThreshold)
	if err != nil {
		return nil, err
	}

	// if STATUS_CHALLENGE_PERIOD is expired, change to STATUS_VERIFIED
	err = k.ChangeToVerifiedFromChallengePeriod(ctx, params.ChallengePeriod)
	if err != nil {
		return nil, err
	}

	// if STATUS_CHALLENGING, tally validity_proofs
	err = k.TallyValidityProofs(ctx, params.ProofPeriod, params.ReplicationFactor)
	if err != nil {
		return nil, err
	}

	// slash epoch moved from PreBlocker
	lastSlashBlockHeight, err := k.LastSlashBlockHeight.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// If LastSlashBlockHeight is not found (e.g., after a new deployment),
			// initialize it with the current block height. This prevents the slash condition
			// from triggering immediately.
			lastSlashBlockHeight = ctx.BlockHeight()
			if err := k.LastSlashBlockHeight.Set(ctx, lastSlashBlockHeight); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	if ctx.BlockHeight() >= lastSlashBlockHeight+int64(params.SlashEpoch) {
		err = k.HandleSlashEpoch(ctx)
		if err != nil {
			return nil, err
		}
		err = k.LastSlashBlockHeight.Set(ctx, ctx.BlockHeight())
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgVerifyDataResponse{}, nil
}
