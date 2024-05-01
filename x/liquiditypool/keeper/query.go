package keeper

import (
	"github.com/sunrise-zone/sunrise-app/x/liquiditypool/types"
)

var _ types.QueryServer = Keeper{}
