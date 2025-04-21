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

	Schema              collections.Schema
	Params              collections.Item[types.Params]
	NextLockupAccountId collections.Map[sdk.AccAddress, uint64]
	LockupAccounts      collections.Map[collections.Pair[[]byte, uint64], types.LockupAccount]

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	feeKeeper     types.FeeKeeper
}

func NewKeeper(
	env appmodule.Environment,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	authority []byte,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	feeKeeper types.FeeKeeper,
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

		Params:              collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		NextLockupAccountId: collections.NewMap(sb, types.NextLockupAccountIdKeyPrefix, "next_lockup_account_id", types.NextLockupAccountIdKeyCodec, collections.Uint64Value),
		LockupAccounts:      collections.NewMap(sb, types.LockupAccountsKeyPrefix, "lockup_accounts", types.LockupAccountsKeyCodec, codec.CollValue[types.LockupAccount](cdc)),

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		feeKeeper:     feeKeeper,
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
