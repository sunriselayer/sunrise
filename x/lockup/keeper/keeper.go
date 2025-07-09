package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
	logger       log.Logger

	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority string

	addressCodec address.Codec

	Schema              collections.Schema
	Params              collections.Item[types.Params]
	NextLockupAccountId collections.Map[sdk.AccAddress, uint64]
	LockupAccounts      collections.Map[collections.Pair[[]byte, uint64], types.LockupAccount]
	Listings            collections.Map[collections.Pair[[]byte, uint64], sdk.Coin]

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	feeKeeper     types.FeeKeeper

	shareclassKeeper types.ShareclassKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	addressCodec address.Codec,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	feeKeeper types.FeeKeeper,
	shareclassKeeper types.ShareclassKeeper,
) Keeper {
	// authority is checked once on startup, so that we don't have to check it again on every message
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		logger:       logger,
		authority:    authority,
		addressCodec: addressCodec,

		Params:              collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		NextLockupAccountId: collections.NewMap(sb, types.NextLockupAccountIdKeyPrefix, "next_lockup_account_id", types.NextLockupAccountIdKeyCodec, collections.Uint64Value),
		LockupAccounts:      collections.NewMap(sb, types.LockupAccountsKeyPrefix, "lockup_accounts", types.LockupAccountsKeyCodec, codec.CollValue[types.LockupAccount](cdc)),
		Listings:            collections.NewMap(sb, types.ListingsKeyPrefix, "listings", types.ListingsKeyCodec, codec.CollValue[sdk.Coin](cdc)),

		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		stakingKeeper:    stakingKeeper,
		feeKeeper:        feeKeeper,
		shareclassKeeper: shareclassKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
