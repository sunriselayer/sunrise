package v0_2_0_test

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
	fromAddress := "sunrise155u042u8wk3al32h3vzxu989jj76k4zcc6d03n"
	toAddresses := []string{
		// new validators
		"sunrise1m63dprapnud2sy3npvw5mgh4nw606u7x5krrhw",
		"sunrise18w30e2qvexwmge4mct99n7mmreczv8sacr322z",
		"sunrise1kw8x5dncdw7ualrx02q4cldcxhsmg5vwtxaxvq",
	}
	// same amount as older validator's one
	coin := sdk.NewInt64Coin("uvrise", 100)

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
		if err := bankkeeper.SendCoins(ctx, fromAddr, toAddr, sdk.NewCoins(coin)); err != nil {
			panic(err)
		}
		ctx.Logger().Info(fmt.Sprintf("send coin [%s] : target [%s]", strconv.Itoa(index), toAddress))

	}
	return nil
}
