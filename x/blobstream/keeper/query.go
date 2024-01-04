package keeper

import (
	"github.com/sunrise-zone/sunrise-app/x/blobstream/types"
)

var _ types.QueryServer = Keeper{}
