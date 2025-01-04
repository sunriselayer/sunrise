package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/tokenconverter/keeper"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func MsgSelfDelegateFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgSelfDelegate] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgSelfDelegate) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgSelfDelegate{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the SelfDelegate simulation

		return []simsx.SimAccount{from}, msg
	}
}
