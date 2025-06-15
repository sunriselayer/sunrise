package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	logger       log.Logger
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	authKeeper    types.AuthKeeper
	bankKeeper    types.BankKeeper
	ybtbaseKeeper types.YbtbaseKeeper

	Schema collections.Schema
	Params collections.Item[types.Params]

	// State
	Tokens              collections.Map[string, types.Token]
	YieldIndex          collections.Map[collections.Pair[string, string], math.LegacyDec]
	UserLastYieldIndex  collections.Map[collections.Triple[string, string, string], math.LegacyDec]
}

func NewKeeper(
	cdc codec.Codec,
	storeService corestore.KVStoreService,
	logger log.Logger,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
	ybtbaseKeeper types.YbtbaseKeeper,
	addressCodec address.Codec,
) (Keeper, error) {
	authority := authtypes.NewModuleAddress(types.GovModuleName)

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:  storeService,
		cdc:           cdc,
		addressCodec:  addressCodec,
		logger:        logger,
		authority:     authority,
		authKeeper:    authKeeper,
		bankKeeper:    bankKeeper,
		ybtbaseKeeper: ybtbaseKeeper,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Tokens: collections.NewMap(sb, types.TokensKey, "tokens", collections.StringKey, codec.CollValue[types.Token](cdc)),
		YieldIndex: collections.NewMap(
			sb,
			types.YieldIndexKey,
			"yield_index",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			sdk.LegacyDecValue,
		),
		UserLastYieldIndex: collections.NewMap(
			sb,
			types.UserLastYieldIndexKey,
			"user_last_yield_index",
			collections.TripleKeyCodec(collections.StringKey, collections.StringKey, collections.StringKey),
			sdk.LegacyDecValue,
		),
	}

	schema, err := sb.Build()
	if err != nil {
		return Keeper{}, err
	}
	k.Schema = schema

	return k, nil
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}

// GetTokenDenom returns the denom for a ybtbrand token
func GetTokenDenom(creator string) string {
	return fmt.Sprintf("ybtbrand/%s", creator)
}

// GetCollateralPoolAddress returns the address of the collateral pool for a token
func GetCollateralPoolAddress(creator string) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("ybtbrand/collateral/%s", creator))
}

// GetYieldPoolAddress returns the address of a yield pool for a token
func GetYieldPoolAddress(creator, yieldDenom string) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("ybtbrand/yield/%s/%s", creator, yieldDenom))
}

// Token operations
func (k Keeper) GetToken(ctx context.Context, creator string) (types.Token, bool) {
	token, err := k.Tokens.Get(ctx, creator)
	if err != nil {
		return types.Token{}, false
	}
	return token, true
}

func (k Keeper) SetToken(ctx context.Context, creator string, token types.Token) error {
	return k.Tokens.Set(ctx, creator, token)
}

func (k Keeper) HasToken(ctx context.Context, creator string) bool {
	has, _ := k.Tokens.Has(ctx, creator)
	return has
}

// YieldIndex operations
func (k Keeper) GetYieldIndex(ctx context.Context, creator, denom string) (math.LegacyDec, bool) {
	index, err := k.YieldIndex.Get(ctx, collections.Join(creator, denom))
	if err != nil {
		return math.LegacyOneDec(), false
	}
	return index, true
}

func (k Keeper) SetYieldIndex(ctx context.Context, creator, denom string, index math.LegacyDec) error {
	return k.YieldIndex.Set(ctx, collections.Join(creator, denom), index)
}

// UserLastYieldIndex operations
func (k Keeper) GetUserLastYieldIndex(ctx context.Context, creator, user, denom string) (math.LegacyDec, bool) {
	index, err := k.UserLastYieldIndex.Get(ctx, collections.Join3(creator, user, denom))
	if err != nil {
		return math.LegacyOneDec(), false
	}
	return index, true
}

func (k Keeper) SetUserLastYieldIndex(ctx context.Context, creator, user, denom string, index math.LegacyDec) error {
	return k.UserLastYieldIndex.Set(ctx, collections.Join3(creator, user, denom), index)
}