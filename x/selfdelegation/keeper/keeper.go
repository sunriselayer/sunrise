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
	LockupAccounts        collections.Map[[]byte, []byte] // lockup account address -> owner address
	SelfDelegationProxies collections.Map[[]byte, []byte] // owner address -> self-delegation proxy address

	accountsKeeper       types.AccountsKeeper
	bankKeeper           types.BankKeeper
	stakingKeeper        types.StakingKeeper
	feeKeeper            types.FeeKeeper
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
	stakingKeeper types.StakingKeeper,
	feeKeeper types.FeeKeeper,
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
		LockupAccounts:        collections.NewMap(sb, types.LockupAccountsKeyPrefix, "lockup_accounts", types.LockupAccountsKeyCodec, collections.BytesValue),
		SelfDelegationProxies: collections.NewMap(sb, types.SelfDelegationProxiesKeyPrefix, "self_delegation_proxies", types.SelfDelegationProxiesKeyCodec, collections.BytesValue),

		accountsKeeper:       accountsKeeper,
		bankKeeper:           bankKeeper,
		stakingKeeper:        stakingKeeper,
		feeKeeper:            feeKeeper,
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
