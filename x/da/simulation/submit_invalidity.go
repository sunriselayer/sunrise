package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/da/keeper"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func MsgSubmitInvalidityFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgSubmitInvalidity] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgSubmitInvalidity) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgSubmitInvalidity{
			Sender: from.AddressBech32,
		}

		// TODO: Handle the SubmitInvalidity simulation

		return []simsx.SimAccount{from}, msg
	}
}
