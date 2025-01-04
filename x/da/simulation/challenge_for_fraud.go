package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"sunrise/x/da/keeper"
	"sunrise/x/da/types"
)

func MsgChallengeForFraudFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgChallengeForFraud] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgChallengeForFraud) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgChallengeForFraud{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the ChallengeForFraud simulation

		return []simsx.SimAccount{from}, msg
	}
}
