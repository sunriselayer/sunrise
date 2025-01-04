package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"sunrise/x/liquiditypool/keeper"
	"sunrise/x/liquiditypool/types"
)

func MsgCreatePoolFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgCreatePool] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgCreatePool) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgCreatePool{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the CreatePool simulation

		return []simsx.SimAccount{from}, msg
	}
}
