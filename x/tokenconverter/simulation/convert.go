package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/tokenconverter/keeper"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func MsgConvertFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgConvert] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgConvert) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgConvert{
			Sender: from.AddressBech32,
		}

		// TODO: Handle the Convert simulation

		return []simsx.SimAccount{from}, msg
	}
}
