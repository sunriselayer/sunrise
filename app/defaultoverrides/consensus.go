package defaultoverrides

import (
	"time"

	cmtcfg "github.com/cometbft/cometbft/config"
	coretypes "github.com/cometbft/cometbft/types"
	"github.com/sunriselayer/sunrise/pkg/appconsts"
)

// DefaultConsensusParams returns a ConsensusParams with a MaxBytes
// determined using a goal square size.
func DefaultConsensusParams() *coretypes.ConsensusParams {
	block := defaultBlockParams()
	evidence := defaultEvidenceParams()
	validatorTmp := coretypes.DefaultValidatorParams()
	validator := coretypes.ValidatorParams{
		PubKeyTypes: validatorTmp.PubKeyTypes,
	}
	return &coretypes.ConsensusParams{
		Block:     block,
		Evidence:  evidence,
		Validator: validator,
		Version: coretypes.VersionParams{
			App: appconsts.LatestVersion,
		},
	}
}

// defaultBlockParams returns a default BlockParams with a MaxBytes determined
// using a goal square size.
func defaultBlockParams() coretypes.BlockParams {
	return coretypes.BlockParams{
		MaxBytes: appconsts.DefaultMaxBytes,
		MaxGas:   -1,
	}
}

// defaultEvidenceParams returns a default EvidenceParams with a MaxAge
// determined using a goal block time.
func defaultEvidenceParams() coretypes.EvidenceParams {
	evdParamsTmp := coretypes.DefaultEvidenceParams()
	evdParams := coretypes.EvidenceParams{
		MaxBytes: evdParamsTmp.MaxBytes,
	}
	evdParams.MaxAgeDuration = appconsts.DefaultUnbondingTime
	evdParams.MaxAgeNumBlocks = int64(appconsts.DefaultUnbondingTime.Seconds())/int64(appconsts.GoalBlockTime.Seconds()) + 1
	return evdParams
}

func DefaultConsensusConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()
	// Set broadcast timeout to be 50 seconds in order to avoid timeouts for long block times
	// TODO: make TimeoutBroadcastTx configurable per https://github.com/celestiaorg/celestia-app/issues/1034
	cfg.RPC.TimeoutBroadcastTxCommit = 50 * time.Second
	cfg.RPC.MaxBodyBytes = int64(8388608) // 8 MiB

	// cfg.Mempool.TTLNumBlocks = 5
	// cfg.Mempool.TTLDuration = time.Duration(cfg.Mempool.TTLNumBlocks) * appconsts.GoalBlockTime
	// Given that there is a stateful transaction size check in CheckTx,
	// We set a loose upper bound on what we expect the transaction to
	// be based on the upper bound size of the entire block for the given
	// version. This acts as a first line of DoS protection
	upperBoundBytes := appconsts.DefaultSquareSizeUpperBound * appconsts.DefaultSquareSizeUpperBound * appconsts.ContinuationSparseShareContentSize
	cfg.Mempool.MaxTxBytes = upperBoundBytes
	// cfg.Mempool.MaxTxsBytes = int64(upperBoundBytes) * cfg.Mempool.TTLNumBlocks
	// cfg.Mempool.Version = "v1" // prioritized mempool

	cfg.Consensus.TimeoutPropose = appconsts.TimeoutPropose
	cfg.Consensus.TimeoutCommit = appconsts.TimeoutCommit
	cfg.Consensus.SkipTimeoutCommit = false

	// cfg.TxIndex.Indexer = "null"
	// cfg.Storage.DiscardABCIResponses = true

	return cfg
}
