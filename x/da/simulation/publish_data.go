package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"sunrise/x/da/keeper"
	"sunrise/x/da/types"
)

func MsgPublishDataFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgPublishData] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgPublishData) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgPublishData{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the PublishData simulation

		return []simsx.SimAccount{from}, msg
	}
}
