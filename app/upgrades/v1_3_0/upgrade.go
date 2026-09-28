// This file contains the upgrade handler for the shutdown upgrade.
package v1_3_0

import (
	"context"
	"fmt"

	"cosmossdk.io/log"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// ValidatorExporter writes the current validator set into the state file.
type ValidatorExporter interface {
	ExportValidators(ctx sdk.Context) ([]cmttypes.GenesisValidator, error)
}

// ConsensusParamsFunc returns the consensus parameters stored at the upgrade height.
type ConsensusParamsFunc func(ctx sdk.Context) cmtproto.ConsensusParams

// CreateUpgradeHandler creates the handler for UpgradeName.
// The proposal that schedules it chooses the block height later.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.JSONCodec,
	homeDir string,
	logger log.Logger,
	bank BankKeeper,
	channels ChannelLister,
	exportValidators ValidatorExporter,
	consensusParams ConsensusParamsFunc,
) upgradetypes.UpgradeHandler {
	return func(goCtx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(goCtx)
		ctx.Logger().Info(fmt.Sprintf("upgrade start: %s at height %d", UpgradeName, ctx.BlockHeight()))

		var validators []cmttypes.GenesisValidator
		if exportValidators != nil {
			exported, err := exportValidators.ExportValidators(ctx)
			if err != nil {
				return nil, fmt.Errorf("CreateUpgradeHandler: export validators: %w", err)
			}
			validators = exported
		}
		var params cmtproto.ConsensusParams
		if consensusParams != nil {
			params = consensusParams(ctx)
		}
		if err := WritePreUpgradeFile(ctx, logger, homeDir, mm, cdc, params, validators); err != nil {
			return nil, err
		}
		if err := SweepIBC(ctx, bank, channels); err != nil {
			return nil, err
		}

		return mm.RunMigrations(goCtx, configurator, fromVM)
	}
}
