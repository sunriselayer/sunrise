package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/liquidstaking/keeper"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func MsgClaimRewardsFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgClaimRewards] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgClaimRewards) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgClaimRewards{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the ClaimRewards simulation

		return []simsx.SimAccount{from}, msg
	}
}
