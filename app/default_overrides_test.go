package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/assert"
	"github.com/sunriselayer/sunrise/app/encoding"
)

// TestDefaultGenesis verifies that the distribution module's genesis state has
// defaults overridden.
func TestDefaultGenesis(t *testing.T) {
	encCfg := encoding.MakeConfig(ModuleEncodingRegisters...)
	dm := distributionModule{}
	raw := dm.DefaultGenesis(encCfg.Codec)
	distributionGenesisState := distributiontypes.GenesisState{}
	encCfg.Codec.MustUnmarshalJSON(raw, &distributionGenesisState)

	// Verify that BaseProposerReward and BonusProposerReward were overridden to 0%.
	assert.Equal(t, sdkmath.LegacyZeroDec(), distributionGenesisState.Params.BaseProposerReward)
	assert.Equal(t, sdkmath.LegacyZeroDec(), distributionGenesisState.Params.BonusProposerReward)

	// Verify that other params weren't overridden.
	assert.Equal(t, distributiontypes.DefaultParams().WithdrawAddrEnabled, distributionGenesisState.Params.WithdrawAddrEnabled)
	assert.Equal(t, distributiontypes.DefaultParams().CommunityTax, distributionGenesisState.Params.CommunityTax)
}

func TestDefaultAppConfig(t *testing.T) {
	cfg := DefaultAppConfig()

	assert.False(t, cfg.API.Enable)
	assert.False(t, cfg.GRPC.Enable)
	assert.False(t, cfg.GRPCWeb.Enable)

	assert.Equal(t, uint64(1500), cfg.StateSync.SnapshotInterval)
	assert.Equal(t, uint32(2), cfg.StateSync.SnapshotKeepRecent)
	assert.Equal(t, "0.002usr", cfg.MinGasPrices)
}
