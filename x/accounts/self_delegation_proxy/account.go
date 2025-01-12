package selfdelegationproxy

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/x/accounts/accountstd"
)

var _ accountstd.Interface = (*SelfDelegationProxyAccount)(nil)

const (
	SELF_DELEGATION_PROXY_ACCOUNT = "self-delegation-proxy"
)

var (
	OwnerPrefix     = collections.NewPrefix(0)
	RootOwnerPrefix = collections.NewPrefix(1)
)

func NewAccount(name string, validatorAddressCodec address.ValidatorAddressCodec) accountstd.AccountCreatorFunc {
	return func(deps accountstd.Dependencies) (string, accountstd.Interface, error) {
		acc := SelfDelegationProxyAccount{
			addressCodec:          deps.AddressCodec,
			validatorAddressCodec: validatorAddressCodec,

			Owner:     collections.NewItem(deps.SchemaBuilder, OwnerPrefix, "owner", collections.BytesValue),
			RootOwner: collections.NewItem(deps.SchemaBuilder, RootOwnerPrefix, "root_owner", collections.BytesValue),
		}

		return name, acc, nil
	}
}

type SelfDelegationProxyAccount struct {
	addressCodec          address.Codec
	validatorAddressCodec address.ValidatorAddressCodec
	// BaseAccount or LockupAccount
	Owner collections.Item[[]byte]
	// BaseAccount
	RootOwner collections.Item[[]byte]
}

func (a SelfDelegationProxyAccount) RegisterInitHandler(builder *accountstd.InitBuilder) {
	accountstd.RegisterInitHandler(builder, a.Init)
}

func (a SelfDelegationProxyAccount) RegisterExecuteHandlers(builder *accountstd.ExecuteBuilder) {
	accountstd.RegisterExecuteHandler(builder, a.Undelegate)
	accountstd.RegisterExecuteHandler(builder, a.CancelUnbonding)
	accountstd.RegisterExecuteHandler(builder, a.WithdrawReward)
	accountstd.RegisterExecuteHandler(builder, a.Send)
}

func (a SelfDelegationProxyAccount) RegisterQueryHandlers(builder *accountstd.QueryBuilder) {
}
