package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"sunrise/x/liquidstaking/keeper"
	"sunrise/x/liquidstaking/types"
)

func SimulateMsgLiquidUnstake(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgLiquidUnstake{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the LiquidUnstake simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "LiquidUnstake simulation not implemented"), nil, nil
	}
}
