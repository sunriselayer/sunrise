package defaultoverrides

import (
	cmtcfg "github.com/cometbft/cometbft/config"
	coretypes "github.com/cometbft/cometbft/types"
)

func DefaultConsensusParams() *coretypes.ConsensusParams {
	return coretypes.DefaultConsensusParams()
}

func DefaultConsensusConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	return cfg
}
