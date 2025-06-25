package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

// DepositMarketIncentive handles MsgDepositMarketIncentive.
// It allows any account to deposit funds into the incentive pool.
func (k msgServer) DepositMarketIncentive(ctx context.Context, msg *types.MsgDepositMarketIncentive) (*types.MsgDepositMarketIncentiveResponse, error) {
	depositor, err := k.accountKeeper.AddressCodec().StringToBytes(msg.Sender)
	if err != nil {
		return nil, err
	}

	amount := sdk.NewCoins(msg.Amount)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.IncentivePoolName, amount); err != nil {
		return nil, err
	}

	return &types.MsgDepositMarketIncentiveResponse{}, nil
}
