package testnode

import (
	"context"
	"fmt"

	"github.com/sunriselayer/sunrise/app/encoding"

	"github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	testencoding "github.com/sunriselayer/sunrise/test/util/encoding"
)

func ReadRecentBlocks(ctx context.Context, rpcAddress string, blocks int64) ([]*types.Block, error) {
	client, err := http.New(rpcAddress, "/websocket")
	if err != nil {
		return nil, err
	}
	status, err := client.Status(ctx)
	if err != nil {
		return nil, err
	}
	if status.SyncInfo.LatestBlockHeight < blocks {
		return nil, fmt.Errorf("latest block height %d is less than requested blocks %d", status.SyncInfo.LatestBlockHeight, blocks)
	}
	return ReadBlockHeights(ctx, rpcAddress, status.SyncInfo.LatestBlockHeight-blocks+1, status.SyncInfo.LatestBlockHeight)
}

func ReadBlockchain(ctx context.Context, rpcAddress string) ([]*types.Block, error) {
	client, err := http.New(rpcAddress, "/websocket")
	if err != nil {
		return nil, err
	}
	status, err := client.Status(ctx)
	if err != nil {
		return nil, err
	}
	return ReadBlockHeights(ctx, rpcAddress, 1, status.SyncInfo.LatestBlockHeight)
}

func ReadBlockHeights(ctx context.Context, rpcAddress string, fromHeight, toHeight int64) ([]*types.Block, error) {
	client, err := http.New(rpcAddress, "/websocket")
	if err != nil {
		return nil, err
	}
	blocks := make([]*types.Block, toHeight-fromHeight+1)
	for i := fromHeight; i <= toHeight; i++ {
		resp, err := client.Block(ctx, &i)
		if err != nil {
			return nil, err
		}
		blocks[i-fromHeight] = resp.Block
	}
	return blocks, nil
}

func DecodeBlockData(data types.Data) ([]types.Tx, error) {
	txHashes := make([]types.Tx, 0)
	for _, txHash := range data.Txs {
		txHashes = append(txHashes, txHash)
	}
	return txHashes, nil
}

func CalculateMeanGasFromRecentBlocks(ctx context.Context, rpcAddress, msgType string, blocks int64) (float64, int64, error) {
	client, err := http.New(rpcAddress, "/websocket")
	if err != nil {
		return 0.0, 0, err
	}
	status, err := client.Status(ctx)
	if err != nil {
		return 0.0, 0, err
	}
	if status.SyncInfo.LatestBlockHeight <= blocks {
		return 0.0, 0, fmt.Errorf("latest block height %d is less than %d", status.SyncInfo.LatestBlockHeight, blocks)
	}
	return CalculateMeanGas(ctx, rpcAddress, msgType, status.SyncInfo.LatestBlockHeight-blocks+1, status.SyncInfo.LatestBlockHeight)
}

func CalculateMeanGas(ctx context.Context, rpcAddress, msgType string, fromHeight int64, toHeight int64) (float64, int64, error) {
	var (
		encCfg   = encoding.MakeConfig(testencoding.ModuleEncodingRegisters...)
		decoder  = encCfg.TxConfig.TxDecoder()
		totalGas int64
		count    int64
		average  = func() float64 {
			if count == 0 {
				return 0
			}
			return float64(totalGas) / float64(count)
		}
	)
	client, err := http.New(rpcAddress, "/websocket")
	if err != nil {
		return 0.0, 0, err
	}

	for height := fromHeight; height <= toHeight; height++ {
		resp, err := client.Block(ctx, &height)
		if err != nil {
			return average(), count, err
		}
		indices := make([]int, 0, len(resp.Block.Data.Txs))
		for i, rawTx := range resp.Block.Data.Txs {
			tx, err := decoder(rawTx)
			if err != nil {
				return average(), count, fmt.Errorf("decoding tx (height: %d): %w", height, err)
			}
			msgs := tx.GetMsgs()
			// multi message transactions are not included
			if len(msgs) == 1 && sdk.MsgTypeURL(msgs[0]) == msgType {
				indices = append(indices, i)
			}
		}
		if len(indices) > 0 {
			results, err := client.BlockResults(ctx, &height)
			if err != nil {
				return average(), count, fmt.Errorf("getting block results (height %d): %w", height, err)
			}
			for _, i := range indices {
				totalGas += results.TxsResults[i].GasUsed
				count++
			}
		}
	}
	return average(), count, nil
}
