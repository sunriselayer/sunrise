package keeper

import (
	"sunrise/x/blobgrant/types"
)

var _ types.QueryServer = Keeper{}
