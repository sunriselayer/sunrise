package keeper

import (
	"github.com/sunrise-zone/sunrise-app/x/blob/types"
)

var _ types.QueryServer = Keeper{}
