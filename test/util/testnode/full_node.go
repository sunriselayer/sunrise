package testnode

import (
	"os"
	"path/filepath"

	"cosmossdk.io/log"
	cmtcfg "github.com/cometbft/cometbft/config"
	tmlog "github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	srvtypes "github.com/cosmos/cosmos-sdk/server/types"
)

// NewCometNode creates a ready to use comet node that operates a single
// validator celestia-app network. It expects that all configuration files are
// already initialized and saved to the baseDir.
func NewCometNode(baseDir string, cfg *UniversalTestingConfig) (*node.Node, srvtypes.Application, error) {
	logger := newLogger(cfg)
	tmLogger := newTMLogger(cfg)
	dbPath := filepath.Join(cfg.TmConfig.RootDir, "data")
	db, err := dbm.NewGoLevelDB("application", dbPath, dbm.OptionsMap{})
	if err != nil {
		return nil, nil, err
	}

	cfg.AppOptions.Set(flags.FlagHome, baseDir)

	app := cfg.AppCreator(logger, db, nil, cfg.AppOptions)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.TmConfig.NodeKeyFile())
	if err != nil {
		return nil, nil, err
	}
	// abciWrapper := server.NewCometABCIWrapper(app)

	tmNode, err := node.NewNode(
		cfg.TmConfig,
		privval.LoadOrGenFilePV(cfg.TmConfig.PrivValidatorKeyFile(), cfg.TmConfig.PrivValidatorStateFile()),
		nodeKey,
		proxy.DefaultClientCreator(cfg.TmConfig.ProxyApp, cfg.TmConfig.ABCI, cfg.TmConfig.DBDir()),
		node.DefaultGenesisDocProviderFunc(cfg.TmConfig),
		cmtcfg.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.TmConfig.Instrumentation),
		tmLogger,
	)

	return tmNode, app, err
}

func newLogger(cfg *UniversalTestingConfig) log.Logger {
	if cfg.SuppressLogs {
		return log.NewNopLogger()
	}
	logger := log.NewLogger(tmlog.NewSyncWriter(os.Stdout))
	return logger
}

func newTMLogger(cfg *UniversalTestingConfig) tmlog.Logger {
	if cfg.SuppressLogs {
		return tmlog.NewNopLogger()
	}
	logger := tmlog.NewTMLogger(tmlog.NewSyncWriter(os.Stdout))
	logger = tmlog.NewFilter(logger, tmlog.AllowError())
	return logger
}
