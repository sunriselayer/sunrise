package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

type Keeper struct {
	appmodule.Environment

	cdc          codec.BinaryCodec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema         collections.Schema
	Params         collections.Item[types.Params]
	LockupAccounts collections.Map[sdk.AccAddress, types.LockupAccount]

	accountKeeper types.AccountKeeper
}

func NewKeeper(
	env appmodule.Environment,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	authority []byte,
	accountKeeper types.AccountKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(env.KVStoreService)

	k := Keeper{
		Environment:  env,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params:         collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		LockupAccounts: collections.NewMap(sb, types.LockupAccountsKeyPrefix, "lockup_accounts", types.LockupAccountsKeyCodec, codec.CollValue[types.LockupAccount](cdc)),

		accountKeeper: accountKeeper,
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
