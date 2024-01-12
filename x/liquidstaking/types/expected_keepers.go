package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here

	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAddressAndPermissions(moduleName string) (addr sdk.AccAddress, permissions []string)
	GetModuleAccountAndPermissions(ctx context.Context, moduleName string) (sdk.ModuleAccountI, []string)
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	SetModuleAccount(ctx context.Context, macc sdk.ModuleAccountI)
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here

	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error

	GetSupply(ctx context.Context, denom string) sdk.Coin
	HasSupply(ctx context.Context, denom string) bool
	// GetPaginatedTotalSupply(ctx context.Context, pagination *query.PageRequest) (sdk.Coins, *query.PageResponse, error)
	IterateTotalSupply(ctx context.Context, cb func(sdk.Coin) bool)
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	HasDenomMetaData(ctx context.Context, denom string) bool
	SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata)
	GetAllDenomMetaData(ctx context.Context) []banktypes.Metadata
	IterateAllDenomMetaData(ctx context.Context, cb func(banktypes.Metadata) bool)

	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	DelegateCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error

	DelegateCoins(ctx context.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error
	UndelegateCoins(ctx context.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error
}

// StakingKeeper expected staking keeper
type StakingKeeper interface {
	ValidatorAddressCodec() address.Codec
	ConsensusAddressCodec() address.Codec
	// iterate through validators by operator address, execute func for each validator
	IterateValidators(context.Context,
		func(index int64, validator stakingtypes.ValidatorI) (stop bool)) error

	Validator(context.Context, sdk.ValAddress) (stakingtypes.ValidatorI, error)            // get a particular validator by operator address
	ValidatorByConsAddr(context.Context, sdk.ConsAddress) (stakingtypes.ValidatorI, error) // get a particular validator by consensus address

	// slash the validator and delegators of the validator, specifying offense height, offense power, and slash fraction
	Slash(context.Context, sdk.ConsAddress, int64, int64, math.LegacyDec) (math.Int, error)
	SlashWithInfractionReason(context.Context, sdk.ConsAddress, int64, int64, math.LegacyDec, stakingtypes.Infraction) (math.Int, error)
	Jail(context.Context, sdk.ConsAddress) error   // jail a validator
	Unjail(context.Context, sdk.ConsAddress) error // unjail a validator

	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(context.Context, sdk.AccAddress, sdk.ValAddress) (stakingtypes.DelegationI, error)
	GetAllValidators(ctx context.Context) ([]stakingtypes.Validator, error)

	// MaxValidators returns the maximum amount of bonded validators
	MaxValidators(context.Context) (uint32, error)

	// IsValidatorJailed returns if the validator is jailed.
	IsValidatorJailed(ctx context.Context, addr sdk.ConsAddress) (bool, error)

	BondDenom(ctx context.Context) (string, error)
	GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error)
	HasReceivingRedelegation(ctx context.Context, delAddr sdk.AccAddress, valDstAddr sdk.ValAddress) (bool, error)
	IterateDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress, cb func(delegation stakingtypes.Delegation) (stop bool)) error
	Delegate(
		ctx context.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus,
		validator stakingtypes.Validator, subtractAccount bool,
	) (newShares math.LegacyDec, err error)
	Unbond(
		ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec,
	) (amount math.Int, err error)
	ValidateUnbondAmount(
		ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amt math.Int,
	) (shares math.LegacyDec, err error)
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
}

type DistributionKeeper interface {
	GetDelegatorWithdrawAddr(ctx context.Context, delAddr sdk.AccAddress) (sdk.AccAddress, error)
	WithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
}
