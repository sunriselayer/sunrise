package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"sunrise/x/liquiditypool/keeper"
	"sunrise/x/liquiditypool/types"
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
