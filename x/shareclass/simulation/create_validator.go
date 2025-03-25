package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/shareclass/keeper"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func MsgCreateValidatorFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgCreateValidator] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgCreateValidator) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgCreateValidator{
			ValidatorAddress: from.AddressBech32,
		}

		// TODO: Handle the CreateValidator simulation

		return []simsx.SimAccount{from}, msg
	}
}
