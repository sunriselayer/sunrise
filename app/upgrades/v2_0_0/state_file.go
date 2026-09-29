package v2_0_0

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"cosmossdk.io/log"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

// WritePreUpgradeFile writes the application state at the start of the upgrade block.
// A disk error is logged and ignored so nodes with different free space still apply the same state transition.
// An export error is returned because it would happen on every node.
func WritePreUpgradeFile(
	ctx sdk.Context,
	logger log.Logger,
	homeDir string,
	mm *module.Manager,
	cdc codec.JSONCodec,
	consensusParams cmtproto.ConsensusParams,
	validators []cmttypes.GenesisValidator,
) error {
	if logger == nil {
		logger = ctx.Logger()
	}
	if homeDir == "" {
		logger.Error("skipped pre-upgrade state file because the node home directory is empty")
		return nil
	}

	appState, err := exportAppState(ctx, mm, cdc)
	if err != nil {
		return err
	}

	doc := &genutiltypes.AppGenesis{
		ChainID:       ctx.ChainID(),
		InitialHeight: ctx.BlockHeight(),
		AppState:      appState,
		Consensus:     genutiltypes.NewConsensusGenesis(consensusParams, validators),
	}
	doc.AppName = "sunrise"
	doc.AppVersion = "v2.0.0"

	dir := filepath.Join(homeDir, "shutdown")
	path := filepath.Join(dir, PreUpgradeFileName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logger.Error("skipped pre-upgrade state file because the directory could not be created", "path", dir, "error", err.Error())
		return nil
	}

	encoded, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("WritePreUpgradeFile: marshal state document at height %d: %w", ctx.BlockHeight(), err)
	}
	sum := sha256.Sum256(encoded)
	sumHex := hex.EncodeToString(sum[:])
	if err := os.WriteFile(path, encoded, 0o600); err != nil {
		logger.Error("skipped pre-upgrade state file because the write failed", "path", path, "sha256", sumHex, "error", err.Error())
		return nil
	}
	logger.Info("wrote pre-upgrade state file", "path", path, "height", ctx.BlockHeight(), "sha256", sumHex)
	return nil
}

func exportAppState(ctx sdk.Context, mm *module.Manager, cdc codec.JSONCodec) (json.RawMessage, error) {
	if mm == nil {
		return nil, fmt.Errorf("exportAppState: module manager is nil")
	}
	if cdc == nil {
		return nil, fmt.Errorf("exportAppState: codec is nil")
	}
	names := mm.ModuleNames()
	sort.Strings(names)
	genState, err := mm.ExportGenesisForModules(ctx, cdc, names)
	if err != nil {
		return nil, fmt.Errorf("exportAppState: modules %v: %w", names, err)
	}
	encoded, err := json.Marshal(genState)
	if err != nil {
		return nil, fmt.Errorf("exportAppState: marshal app state: %w", err)
	}
	return encoded, nil
}
