package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/lockup/keeper"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func MsgClaimRewardsFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgClaimRewards] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgClaimRewards) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgClaimRewards{
			Owner: from.AddressBech32,
		}

		// TODO: Handle the ClaimRewards simulation

		return []simsx.SimAccount{from}, msg
	}
}
