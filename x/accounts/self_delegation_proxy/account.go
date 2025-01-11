package selfdelegationproxy

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/x/accounts/accountstd"
)

var _ accountstd.Interface = (*SelfDelegationProxyAccount)(nil)

var (
	ParentPrefix = collections.NewPrefix(0)
)

func NewAccount(name string) accountstd.AccountCreatorFunc {
	return func(deps accountstd.Dependencies) (string, accountstd.Interface, error) {
		acc := SelfDelegationProxyAccount{
			Parent: collections.NewItem(deps.SchemaBuilder, ParentPrefix, "parent", collections.BytesValue),
		}

		return name, acc, nil
	}
}

type SelfDelegationProxyAccount struct {
	// BaseAccount
	Owner collections.Item[[]byte]
	// BaseAccount or LockupAccount
	Parent collections.Item[[]byte]
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
