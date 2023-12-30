package keeper

import (
	"sunrise/x/blob/types"
)

var _ types.QueryServer = Keeper{}
