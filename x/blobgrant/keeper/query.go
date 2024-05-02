package keeper

import (
	"github.com/sunriselayer/sunrise/x/blobgrant/types"
)

var _ types.QueryServer = Keeper{}
