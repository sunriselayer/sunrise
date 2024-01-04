package keeper

import (
	"github.com/sunrise-zone/sunrise-app/x/blobgrant/types"
)

var _ types.QueryServer = Keeper{}
