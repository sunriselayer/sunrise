package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/lockup/keeper"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func MsgSendFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgSend] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgSend) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgSend{
			Owner: from.AddressBech32,
		}

		// TODO: Handle the Send simulation

		return []simsx.SimAccount{from}, msg
	}
}
