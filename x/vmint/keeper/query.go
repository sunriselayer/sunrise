package keeper

import (
	"github.com/sunriselayer/sunrise/x/vmint/types"
)

var _ types.QueryServer = Keeper{}
