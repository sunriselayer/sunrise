package keeper

import (
	"github.com/sunriselayer/sunrise/x/blob/types"
)

var _ types.QueryServer = Keeper{}
