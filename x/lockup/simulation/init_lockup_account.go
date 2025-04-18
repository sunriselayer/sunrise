package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/lockup/keeper"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func MsgInitLockupAccountFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgInitLockupAccount] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgInitLockupAccount) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgInitLockupAccount{
			Sender: from.AddressBech32,
		}

		// TODO: Handle the InitLockupAccount simulation

		return []simsx.SimAccount{from}, msg
	}
}
