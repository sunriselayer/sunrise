package testnode

import (
	"os"
	"path/filepath"

	"cosmossdk.io/log"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	"github.com/cosmos/cosmos-sdk/client/flags"
	srvtypes "github.com/cosmos/cosmos-sdk/server/types"
)

// NewCometNode creates a ready to use comet node that operates a single
// validator celestia-app network. It expects that all configuration files are
// already initialized and saved to the baseDir.
func NewCometNode(baseDir string, cfg *UniversalTestingConfig) (*node.Node, srvtypes.Application, error) {
	logger := newLogger(cfg)
	dbPath := filepath.Join(cfg.TmConfig.RootDir, "data")
	db, err := dbm.NewGoLevelDB("application", dbPath)
	if err != nil {
		return nil, nil, err
	}

	cfg.AppOptions.Set(flags.FlagHome, baseDir)

	app := cfg.AppCreator(logger, db, nil, cfg.AppOptions)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.TmConfig.NodeKeyFile())
	if err != nil {
		return nil, nil, err
	}

	tmNode, err := node.NewNode(
		cfg.TmConfig,
		privval.LoadOrGenFilePV(cfg.TmConfig.PrivValidatorKeyFile(), cfg.TmConfig.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg.TmConfig),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.TmConfig.Instrumentation),
		logger,
	)

	return tmNode, app, err
}

func newLogger(cfg *UniversalTestingConfig) log.Logger {
	if cfg.SuppressLogs {
		return log.NewNopLogger()
	}
	logger := log.NewLogger(log.NewFilterWriter(os.Stdout, log.AllowError()))
	return logger
}
