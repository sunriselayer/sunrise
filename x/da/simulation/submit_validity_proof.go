package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/da/keeper"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func MsgSubmitValidityProofFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgSubmitValidityProof] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgSubmitValidityProof) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgSubmitValidityProof{
			Sender: from.AddressBech32,
		}

		// TODO: Handle the SubmitValidityProof simulation

		return []simsx.SimAccount{from}, msg
	}
}
