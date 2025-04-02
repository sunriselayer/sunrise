package types

import (
	addresscodec "cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	AddressCodec() addresscodec.Codec
	GetModuleAddress(name string) sdk.AccAddress
}
