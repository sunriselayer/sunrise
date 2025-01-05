package custom

import (
	"fmt"

	serverconfig "github.com/cosmos/cosmos-sdk/server/config"

	"github.com/sunriselayer/sunrise/app/consts"
)

func DefaultServerConfig() *serverconfig.Config {
	cfg := serverconfig.DefaultConfig()

	cfg.MinGasPrices = fmt.Sprintf("%v%s", consts.DefaultMinGasPrice, consts.FeeDenom)

	return cfg
}
