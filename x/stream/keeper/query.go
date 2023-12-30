package keeper

import (
	"sunrise/x/stream/types"
)

var _ types.QueryServer = Keeper{}
