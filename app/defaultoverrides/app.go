package defaultoverrides

import (
	"fmt"

	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/sunriselayer/sunrise/pkg/appconsts"
)

func DefaultAppConfig() *serverconfig.Config {
	cfg := serverconfig.DefaultConfig()
	cfg.API.Enable = false
	cfg.GRPC.Enable = false
	cfg.GRPCWeb.Enable = false

	// the default snapshot interval was determined by picking a large enough
	// value as to not dramatically increase resource requirements while also
	// being greater than zero so that there are more nodes that will serve
	// snapshots to nodes that state sync
	cfg.StateSync.SnapshotInterval = 1500
	cfg.StateSync.SnapshotKeepRecent = 2
	cfg.MinGasPrices = fmt.Sprintf("%v%s", appconsts.DefaultMinGasPrice, appconsts.BondDenom)
	return cfg
}
