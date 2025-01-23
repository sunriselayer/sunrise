package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/sunriselayer/sunrise/x/selfdelegation/types"
)

type Keeper struct {
	appmodule.Environment

	cdc                   codec.BinaryCodec
	addressCodec          address.Codec
	validatorAddressCodec address.ValidatorAddressCodec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema                collections.Schema
	Params                collections.Item[types.Params]
	SelfDelegationProxies collections.Map[[]byte, []byte]

	accountsKeeper       types.AccountsKeeper
	bankKeeper           types.BankKeeper
	tokenConverterKeeper types.TokenConverterKeeper
}

func NewKeeper(
	env appmodule.Environment,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	validatorAddressCodec address.ValidatorAddressCodec,
	authority []byte,
	accountsKeeper types.AccountsKeeper,
	bankKeeper types.BankKeeper,
	tokenConverterKeeper types.TokenConverterKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(env.KVStoreService)

	k := Keeper{
		Environment:           env,
		cdc:                   cdc,
		addressCodec:          addressCodec,
		validatorAddressCodec: validatorAddressCodec,
		authority:             authority,

		Params:                collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		SelfDelegationProxies: collections.NewMap(sb, types.SelfDelegationProxiesKeyPrefix, "self_delegation_proxies", collections.BytesKey, collections.BytesValue),

		accountsKeeper:       accountsKeeper,
		bankKeeper:           bankKeeper,
		tokenConverterKeeper: tokenConverterKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
