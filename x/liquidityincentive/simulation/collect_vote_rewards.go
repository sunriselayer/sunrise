package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func MsgCollectVoteRewardsFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgCollectVoteRewards] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgCollectVoteRewards) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgCollectVoteRewards{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the CollectVoteRewards simulation

		return []simsx.SimAccount{from}, msg
	}
}
