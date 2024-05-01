package keeper

import (
	"github.com/sunriselayer/sunrise-app/x/blobstream/types"
)

var _ types.QueryServer = Keeper{}
