package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func MsgCreatePositionFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgCreatePosition] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgCreatePosition) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgCreatePosition{
			Sender: from.AddressBech32,
		}

		// TODO: Handle the CreatePosition simulation

		return []simsx.SimAccount{from}, msg
	}
}
