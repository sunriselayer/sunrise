package blobstream

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise-app/x/blobstream/keeper"
	"github.com/sunriselayer/sunrise-app/x/blobstream/types"
)

const (
	// InitialLatestAttestationNonce the initial value set in genesis of the latest attestation
	// nonce value in store.
	InitialLatestAttestationNonce = uint64(0)
	// InitialEarliestAvailableAttestationNonce the initial value set in genesis of the earliest
	/// available attestation nonce in store.
	InitialEarliestAvailableAttestationNonce = uint64(1)
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetLatestAttestationNonce(ctx, InitialLatestAttestationNonce)
	// The reason we're setting the earliest available nonce to 1 is because at
	// chain startup, a new valset will always be created. Also, it's easier to
	// set it once here rather than conditionally setting it in abci.EndBlocker
	// which is executed on every block.
	k.SetEarliestAvailableAttestationNonce(ctx, InitialEarliestAvailableAttestationNonce)
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
