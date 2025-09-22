package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/tokenfactory/types"
)

// GetAuthorityMetadata returns the authority metadata for a specific denom
func (k Keeper) GetAuthorityMetadata(ctx sdk.Context, denom string) (types.DenomAuthorityMetadata, error) {
	return k.AuthorityMetadata.Get(ctx, denom)
}

// setAuthorityMetadata stores authority metadata for a specific denom
func (k Keeper) setAuthorityMetadata(ctx sdk.Context, denom string, metadata types.DenomAuthorityMetadata) error {
	err := metadata.Validate()
	if err != nil {
		return err
	}
	return k.AuthorityMetadata.Set(ctx, denom, metadata)
}

func (k Keeper) setAdmin(ctx sdk.Context, denom string, admin string) error {
	metadata, err := k.GetAuthorityMetadata(ctx, denom)
	if err != nil {
		return err
	}

	metadata.Admin = admin

	return k.setAuthorityMetadata(ctx, denom, metadata)
}
