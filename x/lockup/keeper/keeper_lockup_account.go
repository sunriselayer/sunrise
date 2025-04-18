package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k Keeper) LockupAccountAddress(owner string) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.LockupAccountModule(owner))
}
