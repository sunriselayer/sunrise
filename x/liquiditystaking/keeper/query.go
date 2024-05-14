package keeper

import (
	"github.com/sunriselayer/sunrise/x/liquiditystaking/types"
)

var _ types.QueryServer = Keeper{}
