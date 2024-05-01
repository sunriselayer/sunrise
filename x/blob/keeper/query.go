package keeper

import (
	"github.com/sunriselayer/sunrise-app/x/blob/types"
)

var _ types.QueryServer = Keeper{}
