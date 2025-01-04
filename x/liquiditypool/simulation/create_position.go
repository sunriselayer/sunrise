package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"sunrise/x/liquiditypool/keeper"
	"sunrise/x/liquiditypool/types"
)

func MsgCreatePositionFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgCreatePosition] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgCreatePosition) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgCreatePosition{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the CreatePosition simulation

		return []simsx.SimAccount{from}, msg
	}
}
