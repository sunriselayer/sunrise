package v0_2_3_test

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func upgradeSendCoin(
	ctx sdk.Context,
	bankkeeper bankkeeper.Keeper,
) error {
	fromAddress := "sunrise1kw8x5dncdw7ualrx02q4cldcxhsmg5vwtxaxvq" // core dev-2
	toAddresses := []string{
		// new validators
		"sunrise1jy3z69tnk38ar2j2078grqxg2pthk7gpl4zpjd", // 01node
	}
	// same amount as older validator's one
	govCoin := sdk.NewInt64Coin("uvrise", 9000000000000)
	feeCoin := sdk.NewInt64Coin("urise", 10000000)

	fromAddr, err := sdk.AccAddressFromBech32(fromAddress)
	if err != nil {
		panic(err)
	}

	for index, toAddress := range toAddresses {
		toAddr, err := sdk.AccAddressFromBech32(toAddress)
		if err != nil {
			panic(err)
		}
		// if the account is not existent, this method creates account internally
		if err := bankkeeper.SendCoins(ctx, fromAddr, toAddr, sdk.NewCoins(govCoin, feeCoin)); err != nil {
			panic(err)
		}
		ctx.Logger().Info(fmt.Sprintf("send coins [%s] : target [%s]", strconv.Itoa(index), toAddress))

	}
	return nil
}
