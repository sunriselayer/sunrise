package keeper

import (
	"sunrise/x/liquidstaking/types"
)

var _ types.QueryServer = Keeper{}
