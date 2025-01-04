package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func MsgVoteGaugeFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgVoteGauge] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgVoteGauge) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgVoteGauge{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the VoteGauge simulation

		return []simsx.SimAccount{from}, msg
	}
}
