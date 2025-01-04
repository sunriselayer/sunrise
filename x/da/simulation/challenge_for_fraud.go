package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/da/keeper"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func MsgChallengeForFraudFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgChallengeForFraud] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgChallengeForFraud) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgChallengeForFraud{
			Sender: from.AddressBech32,
		}

		// TODO: Handle the ChallengeForFraud simulation

		return []simsx.SimAccount{from}, msg
	}
}
