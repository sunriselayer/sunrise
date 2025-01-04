package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"sunrise/x/da/keeper"
	"sunrise/x/da/types"
)

func MsgSubmitProofFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgSubmitProof] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgSubmitProof) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgSubmitProof{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the SubmitProof simulation

		return []simsx.SimAccount{from}, msg
	}
}
