package keeper

import (
	"github.com/sunrise-zone/sunrise-app/x/sunrise/types"
)

var _ types.QueryServer = Keeper{}
