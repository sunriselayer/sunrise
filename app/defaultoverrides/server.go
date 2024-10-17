package defaultoverrides

import (
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

func DefaultServerConfig() *serverconfig.Config {
	cfg := serverconfig.DefaultConfig()

	cfg.MinGasPrices = "0.002urise"
	return cfg
}
