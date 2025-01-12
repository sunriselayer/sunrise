package selfdelegationproxydepinject

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/depinject"
	"cosmossdk.io/x/accounts/accountstd"

	"github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy"
)

type Inputs struct {
	depinject.In

	ValidatorAddressCodec address.ValidatorAddressCodec
}

func ProvideAccount(in Inputs) accountstd.DepinjectAccount {
	return accountstd.DepinjectAccount{
		MakeAccount: selfdelegationproxy.NewAccount(selfdelegationproxy.SELF_DELEGATION_PROXY_ACCOUNT, in.ValidatorAddressCodec),
	}
}
