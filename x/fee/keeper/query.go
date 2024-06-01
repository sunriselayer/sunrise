package keeper

import (
	"github.com/sunriselayer/sunrise/x/fee/types"
)

var _ types.QueryServer = Keeper{}
