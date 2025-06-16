package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/sunriselayer/sunrise/x/ybtbase/keeper"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

// SimulateMsgSend simulates sending base YBT tokens
func SimulateMsgSend(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgSend{
			FromAddress: simAccount.Address.String(),
		}

		// TODO: Implement proper simulation logic
		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "Send simulation not implemented"), nil, nil
	}
}