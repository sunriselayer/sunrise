package keeper

import (
	"sunrise/x/blobstream/types"
)

var _ types.QueryServer = Keeper{}
