package keeper

import (
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

var _ types.QueryServer = Keeper{}
