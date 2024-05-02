package keeper

import (
	"github.com/sunriselayer/sunrise/x/blobstream/types"
)

var _ types.QueryServer = Keeper{}
