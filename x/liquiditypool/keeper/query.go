package keeper

import (
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

var _ types.QueryServer = Keeper{}
