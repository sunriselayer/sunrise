package selfdelegationproxydepinject

import (
	"cosmossdk.io/depinject"
	"cosmossdk.io/x/accounts/accountstd"

	"github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy"
)

type Inputs struct {
	depinject.In
}

func ProvideAccount(in Inputs) accountstd.DepinjectAccount {
	return accountstd.DepinjectAccount{MakeAccount: selfdelegationproxy.NewAccount("self-delegation-proxy")}
}
