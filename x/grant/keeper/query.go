package keeper

import (
	"sunrise/x/grant/types"
)

var _ types.QueryServer = Keeper{}
