package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/sunriselayer/sunrise/x/fee/types"
)

func (k msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.GetAuthority() != req.Authority {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
	}

	if err := req.Params.Validate(); err != nil {
		return nil, err
	}

	if req.Params.BurnEnabled {
		if req.Params.FeeDenom != req.Params.BurnDenom {
			pool, found, err := k.liquidityPoolKeeper.GetPool(ctx, req.Params.BurnPoolId)
			if err != nil {
				return nil, err
			}
			if !found {
				return nil, errorsmod.Wrapf(types.ErrInvalidPool, "pool %d not found", req.Params.BurnPoolId)
			}

			if !(pool.DenomBase == req.Params.FeeDenom && pool.DenomQuote == req.Params.BurnDenom) && !(pool.DenomBase == req.Params.BurnDenom && pool.DenomQuote == req.Params.FeeDenom) {
				return nil, errorsmod.Wrapf(types.ErrInvalidPool, "pool %d does not have the correct denoms", req.Params.BurnPoolId)
			}
		}
	}

	if err := k.Params.Set(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
