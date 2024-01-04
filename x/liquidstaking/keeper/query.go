package keeper

import (
	"github.com/sunrise-zone/sunrise-app/x/liquidstaking/types"
)

var _ types.QueryServer = Keeper{}
