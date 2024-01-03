package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"time"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here

	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error

	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

// StakingKeeper expected staking keeper (noalias)
type StakingKeeper interface {
	Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)

	GetAllValidators(ctx sdk.Context) (validators []stakingtypes.Validator)
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator

	GetLastTotalPower(ctx sdk.Context) sdkmath.Int
	GetLastValidatorPower(ctx sdk.Context, valAddr sdk.ValAddress) int64

	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI
	GetDelegation(ctx sdk.Context,
		delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, found bool)
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))
	Delegate(
		ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdkmath.Int, tokenSrc stakingtypes.BondStatus,
		validator stakingtypes.Validator, subtractAccount bool,
	) (newShares sdkmath.LegacyDec, err error)

	BondDenom(ctx sdk.Context) (res string)
	Unbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdkmath.LegacyDec) (amount sdkmath.Int, err error)
	UnbondingTime(ctx sdk.Context) (res time.Duration)
	SetUnbondingDelegationEntry(
		ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
		creationHeight int64, minTime time.Time, balance sdkmath.Int,
	) stakingtypes.UnbondingDelegation
	InsertUBDQueue(ctx sdk.Context, ubd stakingtypes.UnbondingDelegation,
		completionTime time.Time)
	ValidateUnbondAmount(
		ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amt sdkmath.Int,
	) (shares sdkmath.LegacyDec, err error)
	GetAllUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress) []stakingtypes.UnbondingDelegation
	BeginRedelegation(
		ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount sdkmath.LegacyDec,
	) (completionTime time.Time, err error)
	GetAllRedelegations(
		ctx sdk.Context, delegator sdk.AccAddress, srcValAddress, dstValAddress sdk.ValAddress,
	) []stakingtypes.Redelegation
	HasReceivingRedelegation(ctx sdk.Context, delAddr sdk.AccAddress, valDstAddr sdk.ValAddress) bool
	BlockValidatorUpdates(ctx sdk.Context) []abci.ValidatorUpdate
	HasMaxUnbondingDelegationEntries(ctx sdk.Context,
		delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) bool
}

// GovKeeper expected gov keeper (noalias)
type GovKeeper interface {
	Tally(ctx sdk.Context, proposal govtypes.Proposal) (passes bool, burnDeposits bool, tallyResults govtypes.TallyResult)
	GetVotes(ctx sdk.Context, proposalID uint64) (votes govtypes.Votes)
	GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (vote govtypes.Vote, found bool)
	GetProposal(ctx sdk.Context, proposalID uint64) (govtypes.Proposal, bool)
	SetProposal(ctx sdk.Context, proposal govtypes.Proposal)
	GetProposals(ctx sdk.Context) (proposals govtypes.Proposals)
	AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options govtypes.WeightedVoteOptions) error
}

// DistrKeeper expected distribution keeper (noalias)
type DistrKeeper interface {
	IncrementValidatorPeriod(ctx sdk.Context, val stakingtypes.ValidatorI) uint64
	CalculateDelegationRewards(ctx sdk.Context, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdkmath.LegacyDecCoins)
	WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
}

// SlashingKeeper expected slashing keeper (noalias)
type SlashingKeeper interface {
	IsTombstoned(ctx sdk.Context, consAddr sdk.ConsAddress) bool
}

// StakingHooks event hooks for staking validator object (noalias)
type StakingHooks interface {
	AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress)                           // Must be called when a validator is created
	AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) // Must be called when a validator is deleted

	BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress)        // Must be called when a delegation is created
	BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) // Must be called when a delegation's shares are modified
	AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress)
	BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdkmath.LegacyDec)
}
