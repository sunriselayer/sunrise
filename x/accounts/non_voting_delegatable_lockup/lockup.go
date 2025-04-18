// https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-rc.1/x/accounts/defaults/lockup/lockup.go
package lockup

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/cosmos/gogoproto/proto"

	"cosmossdk.io/collections"
	collcodec "cosmossdk.io/collections/codec"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/header"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"cosmossdk.io/x/accounts/accountstd"
	banktypes "cosmossdk.io/x/bank/types"

	distrtypes "cosmossdk.io/x/distribution/types"
	stakingtypes "cosmossdk.io/x/staking/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	// <sunsise>
	lockuptypes "github.com/sunriselayer/sunrise/x/accounts/non_voting_delegatable_lockup/v1"
	feetypes "github.com/sunriselayer/sunrise/x/fee/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
	// </sunrise>
)

var (
	OriginalLockingPrefix  = collections.NewPrefix(0)
	DelegatedFreePrefix    = collections.NewPrefix(1)
	DelegatedLockingPrefix = collections.NewPrefix(2)
	EndTimePrefix          = collections.NewPrefix(3)
	StartTimePrefix        = collections.NewPrefix(4)
	LockingPeriodsPrefix   = collections.NewPrefix(5)
	OwnerPrefix            = collections.NewPrefix(6)
	UnbondEntriesPrefix    = collections.NewPrefix(7)
)

var (
	CONTINUOUS_LOCKING_ACCOUNT = "non-voting-delegatable-continuous-locking-account"
	DELAYED_LOCKING_ACCOUNT    = "non-voting-delegatable-delayed-locking-account"
	PERIODIC_LOCKING_ACCOUNT   = "non-voting-delegatable-periodic-locking-account"
	PERMANENT_LOCKING_ACCOUNT  = "non-voting-delegatable-permanent-locking-account"
)

type getLockedCoinsFunc = func(ctx context.Context, time time.Time, denoms ...string) (sdk.Coins, error)

// newBaseLockup creates a new BaseLockup object.
func newBaseLockup(d accountstd.Dependencies) *BaseLockup {
	BaseLockup := &BaseLockup{
		Owner:            collections.NewItem(d.SchemaBuilder, OwnerPrefix, "owner", collections.BytesValue),
		OriginalLocking:  collections.NewMap(d.SchemaBuilder, OriginalLockingPrefix, "original_locking", collections.StringKey, sdk.IntValue),
		DelegatedFree:    collections.NewMap(d.SchemaBuilder, DelegatedFreePrefix, "delegated_free", collections.StringKey, sdk.IntValue),
		DelegatedLocking: collections.NewMap(d.SchemaBuilder, DelegatedLockingPrefix, "delegated_locking", collections.StringKey, sdk.IntValue),
		UnbondEntries:    collections.NewMap(d.SchemaBuilder, UnbondEntriesPrefix, "unbond_entries", collections.StringKey, codec.CollValue[lockuptypes.UnbondingEntries](d.LegacyStateCodec)),
		addressCodec:     d.AddressCodec,
		headerService:    d.Environment.HeaderService,
		EndTime:          collections.NewItem(d.SchemaBuilder, EndTimePrefix, "end_time", collcodec.KeyToValueCodec[time.Time](sdk.TimeKey)),
	}

	return BaseLockup
}

type BaseLockup struct {
	// Owner is the address of the account owner.
	Owner            collections.Item[[]byte]
	OriginalLocking  collections.Map[string, math.Int]
	DelegatedFree    collections.Map[string, math.Int]
	DelegatedLocking collections.Map[string, math.Int]
	// map val address to unbonding entries
	UnbondEntries collections.Map[string, lockuptypes.UnbondingEntries]
	addressCodec  address.Codec
	headerService header.Service
	// lockup end time.
	EndTime collections.Item[time.Time]
}

func (bva *BaseLockup) Init(ctx context.Context, msg *lockuptypes.MsgInitNonVotingDelegatableLockupAccount) (
	*lockuptypes.MsgInitNonVotingDelegatableLockupAccountResponse, error,
) {
	owner, err := bva.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid 'owner' address: %s", err)
	}
	err = bva.Owner.Set(ctx, owner)
	if err != nil {
		return nil, err
	}

	funds := accountstd.Funds(ctx)

	sortedAmt := funds.Sort()
	for _, coin := range sortedAmt {
		err = bva.OriginalLocking.Set(ctx, coin.Denom, coin.Amount)
		if err != nil {
			return nil, err
		}
	}

	bondDenom, err := getStakingDenom(ctx)
	if err != nil {
		return nil, err
	}

	// Set initial value for all locked token
	err = bva.DelegatedFree.Set(ctx, bondDenom, math.ZeroInt())
	if err != nil {
		return nil, err
	}

	// Set initial value for all locked token
	err = bva.DelegatedLocking.Set(ctx, bondDenom, math.ZeroInt())
	if err != nil {
		return nil, err
	}

	err = bva.EndTime.Set(ctx, msg.EndTime)
	if err != nil {
		return nil, err
	}

	return &lockuptypes.MsgInitNonVotingDelegatableLockupAccountResponse{}, nil
}

// func (bva *BaseLockup) Delegate(
// 	ctx context.Context, msg *lockuptypes.MsgDelegate, getLockedCoinsFunc getLockedCoinsFunc,
// ) (
// 	*lockuptypes.MsgExecuteMessagesResponse, error,
// ) {
// 	err := bva.checkSender(ctx, msg.Sender)
// 	if err != nil {
// 		return nil, err
// 	}
// 	whoami := accountstd.Whoami(ctx)
// 	delegatorAddress, err := bva.addressCodec.BytesToString(whoami)
// 	if err != nil {
// 		return nil, err
// 	}

// 	hs := bva.headerService.HeaderInfo(ctx)

// 	balance, err := bva.getBalance(ctx, delegatorAddress, msg.Amount.Denom)
// 	if err != nil {
// 		return nil, err
// 	}
// 	lockedCoins, err := getLockedCoinsFunc(ctx, hs.Time, msg.Amount.Denom)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// refresh ubd entries to make sure delegation locking amount is up to date
// 	err = bva.checkUnbondingEntriesMature(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = bva.TrackDelegation(
// 		ctx,
// 		sdk.Coins{*balance},
// 		lockedCoins,
// 		sdk.Coins{msg.Amount},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	msgDelegate := &stakingtypes.MsgDelegate{
// 		DelegatorAddress: delegatorAddress,
// 		ValidatorAddress: msg.ValidatorAddress,
// 		Amount:           msg.Amount,
// 	}
// 	resp, err := sendMessage(ctx, msgDelegate)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &lockuptypes.MsgExecuteMessagesResponse{Responses: resp}, nil
// }

// func (bva *BaseLockup) Undelegate(
// 	ctx context.Context, msg *lockuptypes.MsgUndelegate,
// ) (
// 	*lockuptypes.MsgExecuteMessagesResponse, error,
// ) {
// 	err := bva.checkSender(ctx, msg.Sender)
// 	if err != nil {
// 		return nil, err
// 	}
// 	whoami := accountstd.Whoami(ctx)
// 	delegatorAddress, err := bva.addressCodec.BytesToString(whoami)
// 	if err != nil {
// 		return nil, err
// 	}

// 	msgUndelegate := &stakingtypes.MsgUndelegate{
// 		DelegatorAddress: delegatorAddress,
// 		ValidatorAddress: msg.ValidatorAddress,
// 		Amount:           msg.Amount,
// 	}
// 	resp, err := sendMessage(ctx, msgUndelegate)
// 	if err != nil {
// 		return nil, err
// 	}

// 	header := bva.headerService.HeaderInfo(ctx)

// 	msgUndelegateResp, err := accountstd.UnpackAny[stakingtypes.MsgUndelegateResponse](resp[0])
// 	if err != nil {
// 		return nil, err
// 	}

// 	isNewEntry := true
// 	entries, err := bva.UnbondEntries.Get(ctx, msg.ValidatorAddress)
// 	if err != nil {
// 		if errorsmod.IsOf(err, collections.ErrNotFound) {
// 			entries = lockuptypes.UnbondingEntries{
// 				Entries: []*lockuptypes.UnbondingEntry{},
// 			}
// 		} else {
// 			return nil, err
// 		}
// 	}

// 	for i, entry := range entries.Entries {
// 		if entry.CreationHeight == header.Height && entry.EndTime.Equal(msgUndelegateResp.CompletionTime) {
// 			entry.Amount = entry.Amount.Add(msg.Amount)

// 			// update the entry
// 			entries.Entries[i] = entry
// 			isNewEntry = false
// 			break
// 		}
// 	}

// 	if isNewEntry {
// 		entries.Entries = append(entries.Entries, &lockuptypes.UnbondingEntry{
// 			EndTime:          msgUndelegateResp.CompletionTime,
// 			Amount:           msgUndelegateResp.Amount,
// 			ValidatorAddress: msg.ValidatorAddress,
// 			CreationHeight:   header.Height,
// 		})
// 	}

// 	err = bva.UnbondEntries.Set(ctx, msg.ValidatorAddress, entries)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &lockuptypes.MsgExecuteMessagesResponse{Responses: resp}, nil
// }

// func (bva *BaseLockup) WithdrawReward(
// 	ctx context.Context, msg *lockuptypes.MsgWithdrawReward,
// ) (
// 	*lockuptypes.MsgExecuteMessagesResponse, error,
// ) {
// 	err := bva.checkSender(ctx, msg.Sender)
// 	if err != nil {
// 		return nil, err
// 	}
// 	whoami := accountstd.Whoami(ctx)
// 	delegatorAddress, err := bva.addressCodec.BytesToString(whoami)
// 	if err != nil {
// 		return nil, err
// 	}

// 	msgWithdraw := &distrtypes.MsgWithdrawDelegatorReward{
// 		DelegatorAddress: delegatorAddress,
// 		ValidatorAddress: msg.ValidatorAddress,
// 	}
// 	responses, err := sendMessage(ctx, msgWithdraw)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &lockuptypes.MsgExecuteMessagesResponse{Responses: responses}, nil
// }

// <sunrise>
func (bva *BaseLockup) Delegate(
	ctx context.Context, msg *lockuptypes.MsgDelegate, getLockedCoinsFunc getLockedCoinsFunc,
) (
	*lockuptypes.MsgExecuteMessagesResponse, error,
) {
	err := bva.checkSender(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}
	whoami := accountstd.Whoami(ctx)
	delegatorAddress, err := bva.addressCodec.BytesToString(whoami)
	if err != nil {
		return nil, err
	}

	hs := bva.headerService.HeaderInfo(ctx)

	balance, err := bva.getBalance(ctx, delegatorAddress, msg.Amount.Denom)
	if err != nil {
		return nil, err
	}
	lockedCoins, err := getLockedCoinsFunc(ctx, hs.Time, msg.Amount.Denom)
	if err != nil {
		return nil, err
	}

	// refresh ubd entries to make sure delegation locking amount is up to date
	err = bva.checkUnbondingEntriesMature(ctx)
	if err != nil {
		return nil, err
	}

	err = bva.TrackDelegation(
		ctx,
		sdk.Coins{*balance},
		lockedCoins,
		sdk.Coins{msg.Amount},
	)
	if err != nil {
		return nil, err
	}

	msgDelegate := &shareclasstypes.MsgNonVotingDelegate{
		Sender:           delegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           msg.Amount,
	}
	resp, err := sendMessage(ctx, msgDelegate)
	if err != nil {
		return nil, err
	}

	return &lockuptypes.MsgExecuteMessagesResponse{Responses: resp}, nil
}

func (bva *BaseLockup) Undelegate(
	ctx context.Context, msg *lockuptypes.MsgUndelegate,
) (
	*lockuptypes.MsgExecuteMessagesResponse, error,
) {
	err := bva.checkSender(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}
	whoami := accountstd.Whoami(ctx)
	delegatorAddress, err := bva.addressCodec.BytesToString(whoami)
	if err != nil {
		return nil, err
	}

	msgUndelegate := &shareclasstypes.MsgNonVotingUndelegate{
		Sender:           delegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           msg.Amount,
	}
	resp, err := sendMessage(ctx, msgUndelegate)
	if err != nil {
		return nil, err
	}

	header := bva.headerService.HeaderInfo(ctx)

	msgUndelegateResp, err := accountstd.UnpackAny[shareclasstypes.MsgNonVotingUndelegateResponse](resp[0])
	if err != nil {
		return nil, err
	}

	isNewEntry := true
	entries, err := bva.UnbondEntries.Get(ctx, msg.ValidatorAddress)
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			entries = lockuptypes.UnbondingEntries{
				Entries: []*lockuptypes.UnbondingEntry{},
			}
		} else {
			return nil, err
		}
	}

	for i, entry := range entries.Entries {
		if entry.CreationHeight == header.Height && entry.EndTime.Equal(msgUndelegateResp.CompletionTime) {
			entry.Amount = entry.Amount.Add(msgUndelegateResp.Amount) // <sunrise />

			// update the entry
			entries.Entries[i] = entry
			isNewEntry = false
			break
		}
	}

	if isNewEntry {
		entries.Entries = append(entries.Entries, &lockuptypes.UnbondingEntry{
			EndTime:          msgUndelegateResp.CompletionTime,
			Amount:           msgUndelegateResp.Amount,
			ValidatorAddress: msg.ValidatorAddress,
			CreationHeight:   header.Height,
		})
	}

	err = bva.UnbondEntries.Set(ctx, msg.ValidatorAddress, entries)
	if err != nil {
		return nil, err
	}

	return &lockuptypes.MsgExecuteMessagesResponse{Responses: resp}, nil
}

func (bva *BaseLockup) WithdrawReward(
	ctx context.Context, msg *lockuptypes.MsgWithdrawReward,
) (
	*lockuptypes.MsgExecuteMessagesResponse, error,
) {
	err := bva.checkSender(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}
	whoami := accountstd.Whoami(ctx)
	delegatorAddress, err := bva.addressCodec.BytesToString(whoami)
	if err != nil {
		return nil, err
	}

	msgWithdraw := &distrtypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
	}
	responses, err := sendMessage(ctx, msgWithdraw)
	if err != nil {
		return nil, err
	}

	return &lockuptypes.MsgExecuteMessagesResponse{Responses: responses}, nil
}

func (bva *BaseLockup) SendCoins(
	ctx context.Context, msg *lockuptypes.MsgSend, getLockedCoinsFunc getLockedCoinsFunc,
) (
	*lockuptypes.MsgExecuteMessagesResponse, error,
) {
	err := bva.checkSender(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}
	whoami := accountstd.Whoami(ctx)
	fromAddress, err := bva.addressCodec.BytesToString(whoami)
	if err != nil {
		return nil, err
	}

	hs := bva.headerService.HeaderInfo(ctx)

	if err := msg.Amount.Validate(); err != nil {
		return nil, err
	}

	lockedCoins, err := getLockedCoinsFunc(ctx, hs.Time, msg.Amount.Denoms()...)
	if err != nil {
		return nil, err
	}

	err = bva.checkTokensSendable(ctx, fromAddress, msg.Amount, lockedCoins)
	if err != nil {
		return nil, err
	}

	msgSend := &banktypes.MsgSend{
		FromAddress: fromAddress,
		ToAddress:   msg.ToAddress,
		Amount:      msg.Amount,
	}
	resp, err := sendMessage(ctx, msgSend)
	if err != nil {
		return nil, err
	}

	return &lockuptypes.MsgExecuteMessagesResponse{Responses: resp}, nil
}

func (bva *BaseLockup) checkSender(ctx context.Context, sender string) error {
	owner, err := bva.Owner.Get(ctx)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err.Error())
	}
	senderBytes, err := bva.addressCodec.StringToBytes(sender)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err.Error())
	}
	if !bytes.Equal(owner, senderBytes) {
		return errors.New("sender is not the owner of this vesting account")
	}

	return nil
}

func sendMessage(ctx context.Context, msg proto.Message) ([]*codectypes.Any, error) {
	asAny, err := accountstd.PackAny(msg)
	if err != nil {
		return nil, err
	}

	return accountstd.ExecModuleAnys(ctx, []*codectypes.Any{asAny})
}

// <sunrise />
func getStakingDenom(ctx context.Context) (string, error) {
	// // Query account balance for the sent denom
	// resp, err := accountstd.QueryModule[*stakingtypes.QueryParamsResponse](ctx, &stakingtypes.QueryParamsRequest{})
	// if err != nil {
	// 	return "", err
	// }

	// return resp.Params.BondDenom, nil
	resp, err := accountstd.QueryModule[*feetypes.QueryParamsResponse](ctx, &feetypes.QueryParamsRequest{})
	if err != nil {
		return "", err
	}

	return resp.Params.FeeDenom, nil
}

// checkUnbondingEntriesMature iterates through all the unbonding entries and check if any of the entries are matured and handled.
func (bva *BaseLockup) checkUnbondingEntriesMature(ctx context.Context) error {
	whoami := accountstd.Whoami(ctx)
	delAddr, err := bva.addressCodec.BytesToString(whoami)
	if err != nil {
		return err
	}

	currentTime := bva.headerService.HeaderInfo(ctx).Time

	removeKeys := []string{}
	err = bva.UnbondEntries.Walk(ctx, nil, func(key string, value lockuptypes.UnbondingEntries) (stop bool, err error) {
		for i := 0; i < len(value.Entries); i++ {
			entry := value.Entries[i]
			// if not mature then skip
			if entry.EndTime.After(currentTime) {
				return false, nil
			}

			stakingUnbonding, err := bva.getUnbondingEntries(ctx, delAddr, key)
			if err != nil {
				// if ubd delegation is empty then skip the next iteration check
				if !errorsmod.IsOf(err, stakingtypes.ErrNoUnbondingDelegation) {
					return true, err
				}
			}

			found := false
			// check if the entry is still exist in the unbonding entries
			for _, e := range stakingUnbonding {
				if e.CompletionTime.Equal(entry.EndTime) && entry.CreationHeight == entry.CreationHeight {
					found = true
					break
				}
			}

			// if not found or ubd delegation is empty then assume ubd entry is being handled
			if !found {
				err = bva.TrackUndelegation(ctx, sdk.NewCoins(entry.Amount))
				if err != nil {
					return true, err
				}

				// remove entry
				value.Entries = append(value.Entries[:i], value.Entries[i+1:]...)
				i--
			}
		}

		if len(value.Entries) == 0 {
			removeKeys = append(removeKeys, key)
		} else {
			err = bva.UnbondEntries.Set(ctx, key, value)
			if err != nil {
				return true, err
			}
		}

		return false, nil
	})

	for _, key := range removeKeys {
		err = bva.UnbondEntries.Remove(ctx, key)
		if err != nil {
			return err
		}
	}
	return err
}

// TrackDelegation tracks a delegation amount for any given lockup account type
// given the amount of coins currently being locked and the current account balance
// of the delegation denominations.
//
// CONTRACT: The account's coins, delegation coins, locked coins, and delegated
// locking coins must be sorted.
func (bva *BaseLockup) TrackDelegation(
	ctx context.Context, balance, lockedCoins, amount sdk.Coins,
) error {
	bondDenom, err := getStakingDenom(ctx)
	if err != nil {
		return err
	}

	delAmt := amount.AmountOf(bondDenom)
	baseAmt := balance.AmountOf(bondDenom)

	// return error if the delegation amount is zero or if the base coins does not
	// exceed the desired delegation amount.
	if delAmt.IsZero() || baseAmt.LT(delAmt) {
		return sdkerrors.ErrInvalidCoins.Wrap("delegation attempt with zero coins for staking denom or insufficient funds")
	}

	lockedAmt := lockedCoins.AmountOf(bondDenom)
	delLockingAmt, err := bva.DelegatedLocking.Get(ctx, bondDenom)
	if err != nil {
		return err
	}
	delFreeAmt, err := bva.DelegatedFree.Get(ctx, bondDenom)
	if err != nil {
		return err
	}

	// compute x and y per the specification, where:
	// X := min(max(V - DV, 0), D)
	// Y := D - X
	x := math.MinInt(math.MaxInt(lockedAmt.Sub(delLockingAmt), math.ZeroInt()), delAmt)
	y := delAmt.Sub(x)

	delLockingCoin := sdk.NewCoin(bondDenom, delLockingAmt)
	delFreeCoin := sdk.NewCoin(bondDenom, delFreeAmt)
	if !x.IsZero() {
		xCoin := sdk.NewCoin(bondDenom, x)
		newDelLocking := delLockingCoin.Add(xCoin)
		err = bva.DelegatedLocking.Set(ctx, bondDenom, newDelLocking.Amount)
		if err != nil {
			return err
		}
	}

	if !y.IsZero() {
		yCoin := sdk.NewCoin(bondDenom, y)
		newDelFree := delFreeCoin.Add(yCoin)
		err = bva.DelegatedFree.Set(ctx, bondDenom, newDelFree.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}

// TrackUndelegation tracks an undelegation amount by setting the necessary
// values by which delegated locking and delegated free need to decrease and
// by which amount the base coins need to increase.
//
// NOTE: The undelegation (bond refund) amount may exceed the delegated
// locking (bond) amount due to the way undelegation truncates the bond refund,
// which can increase the validator's exchange rate (tokens/shares) slightly if
// the undelegated tokens are non-integral.
//
// CONTRACT: The account's coins and undelegation coins must be sorted.
func (bva *BaseLockup) TrackUndelegation(ctx context.Context, amount sdk.Coins) error {
	bondDenom, err := getStakingDenom(ctx)
	if err != nil {
		return err
	}
	delAmt := amount.AmountOf(bondDenom)
	// return error if the undelegation amount is zero
	if delAmt.IsZero() {
		return sdkerrors.ErrInvalidCoins.Wrap("undelegation attempt with zero coins for staking denom")
	}
	delFreeAmt, err := bva.DelegatedFree.Get(ctx, bondDenom)
	if err != nil {
		return err
	}
	delLockingAmt, err := bva.DelegatedLocking.Get(ctx, bondDenom)
	if err != nil {
		return err
	}

	// compute x and y per the specification, where:
	// X := min(DF, D)
	// Y := min(DV, D - X)
	x := math.MinInt(delFreeAmt, delAmt)
	y := math.MinInt(delLockingAmt, delAmt.Sub(x))

	delLockingCoin := sdk.NewCoin(bondDenom, delLockingAmt)
	delFreeCoin := sdk.NewCoin(bondDenom, delFreeAmt)
	if !x.IsZero() {
		xCoin := sdk.NewCoin(bondDenom, x)
		newDelFree := delFreeCoin.Sub(xCoin)
		err = bva.DelegatedFree.Set(ctx, bondDenom, newDelFree.Amount)
		if err != nil {
			return err
		}
	}

	if !y.IsZero() {
		yCoin := sdk.NewCoin(bondDenom, y)
		newDelLocking := delLockingCoin.Sub(yCoin)
		err = bva.DelegatedLocking.Set(ctx, bondDenom, newDelLocking.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bva BaseLockup) getBalance(ctx context.Context, sender, denom string) (*sdk.Coin, error) {
	// Query account balance for the sent denom
	resp, err := accountstd.QueryModule[*banktypes.QueryBalanceResponse](ctx, &banktypes.QueryBalanceRequest{Address: sender, Denom: denom})
	if err != nil {
		return nil, err
	}

	return resp.Balance, nil
}

func (bva BaseLockup) getUnbondingEntries(ctx context.Context, delAddr, valAddr string) ([]stakingtypes.UnbondingDelegationEntry, error) {
	resp, err := accountstd.QueryModule[*stakingtypes.QueryUnbondingDelegationResponse](
		ctx, &stakingtypes.QueryUnbondingDelegationRequest{DelegatorAddr: delAddr, ValidatorAddr: valAddr},
	)
	if err != nil {
		return nil, err
	}

	return resp.Unbond.Entries, nil
}

func (bva BaseLockup) checkTokensSendable(ctx context.Context, sender string, amount, lockedCoins sdk.Coins) error {
	// Check if any sent tokens is exceeds lockup account balances
	for _, coin := range amount {
		balance, err := bva.getBalance(ctx, sender, coin.Denom)
		if err != nil {
			return err
		}
		lockedAmt := lockedCoins.AmountOf(coin.Denom)

		// get lockedCoin from that are not bonded for the sent denom
		notBondedLockedCoin, err := bva.GetNotBondedLockedCoin(ctx, sdk.NewCoin(coin.Denom, lockedAmt), coin.Denom)
		if err != nil {
			return err
		}

		spendable, hasNeg := sdk.Coins{*balance}.SafeSub(notBondedLockedCoin)
		if hasNeg {
			return errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
				"locked amount exceeds account balance funds: %s > %s", notBondedLockedCoin, balance)
		}

		if _, hasNeg := spendable.SafeSub(coin); hasNeg {
			if len(spendable) == 0 {
				spendable = sdk.Coins{sdk.NewCoin(coin.Denom, math.ZeroInt())}
			}
			return errorsmod.Wrapf(
				sdkerrors.ErrInsufficientFunds,
				"spendable balance %s is smaller than %s",
				spendable, coin,
			)
		}
	}

	return nil
}

// IterateCoinEntries iterates over all the CoinEntries entries.
func (bva BaseLockup) IterateCoinEntries(
	ctx context.Context,
	entries collections.Map[string, math.Int],
	cb func(denom string, value math.Int) (bool, error),
) error {
	err := entries.Walk(ctx, nil, func(key string, value math.Int) (stop bool, err error) {
		return cb(key, value)
	})
	return err
}

// GetNotBondedLockedCoin returns the coin that are not spendable that are not bonded by denom
// for a lockup account. If the coin by the provided denom are not locked, an coin with zero amount is returned.
func (bva BaseLockup) GetNotBondedLockedCoin(ctx context.Context, lockedCoin sdk.Coin, denom string) (sdk.Coin, error) {
	// refresh the unbonding entries
	err := bva.checkUnbondingEntriesMature(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	bondDenom, err := getStakingDenom(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	// if not bond denom then return the full locked coin
	if bondDenom != denom {
		return lockedCoin, nil
	}

	delegatedLockingAmt, err := bva.DelegatedLocking.Get(ctx, denom)
	if err != nil {
		return sdk.Coin{}, err
	}

	x := math.MinInt(lockedCoin.Amount, delegatedLockingAmt)
	lockedAmt := lockedCoin.Amount.Sub(x)

	return sdk.NewCoin(denom, lockedAmt), nil
}

// QueryLockupAccountBaseInfo returns a lockup account's info
func (bva BaseLockup) QueryLockupAccountBaseInfo(ctx context.Context, _ *lockuptypes.QueryLockupAccountInfoRequest) (
	*lockuptypes.QueryLockupAccountInfoResponse, error,
) {
	owner, err := bva.Owner.Get(ctx)
	if err != nil {
		return nil, err
	}

	ownerAddress, err := bva.addressCodec.BytesToString(owner)
	if err != nil {
		return nil, err
	}

	endTime, err := bva.EndTime.Get(ctx)
	if err != nil {
		return nil, err
	}

	originalLocking := sdk.Coins{}
	err = bva.IterateCoinEntries(ctx, bva.OriginalLocking, func(key string, value math.Int) (stop bool, err error) {
		originalLocking = append(originalLocking, sdk.NewCoin(key, value))
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	bondDenom, err := getStakingDenom(ctx)
	if err != nil {
		return nil, err
	}

	delegatedLockingAmt, err := bva.DelegatedLocking.Get(ctx, bondDenom)
	if err != nil {
		return nil, err
	}
	delegatedLocking := sdk.NewCoins(sdk.NewCoin(bondDenom, delegatedLockingAmt))

	delegatedFreeAmt, err := bva.DelegatedFree.Get(ctx, bondDenom)
	if err != nil {
		return nil, err
	}
	delegatedFree := sdk.NewCoins(sdk.NewCoin(bondDenom, delegatedFreeAmt))

	return &lockuptypes.QueryLockupAccountInfoResponse{
		Owner:            ownerAddress,
		OriginalLocking:  originalLocking,
		DelegatedLocking: delegatedLocking,
		DelegatedFree:    delegatedFree,
		EndTime:          &endTime,
	}, nil
}

func (bva BaseLockup) QueryUnbondingEntries(ctx context.Context, req *lockuptypes.QueryUnbondingEntriesRequest) (
	*lockuptypes.QueryUnbondingEntriesResponse, error,
) {
	entries, err := bva.UnbondEntries.Get(ctx, req.ValidatorAddress)
	if err != nil {
		if !errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, err
		}

		entries = lockuptypes.UnbondingEntries{
			Entries: []*lockuptypes.UnbondingEntry{},
		}
	}

	return &lockuptypes.QueryUnbondingEntriesResponse{
		UnbondingEntries: entries.Entries,
	}, nil
}

func (bva BaseLockup) QuerySpendableTokens(ctx context.Context, lockedCoins sdk.Coins) (
	*lockuptypes.QuerySpendableAmountResponse, error,
) {
	whoami := accountstd.Whoami(ctx)
	accAddr, err := bva.addressCodec.BytesToString(whoami)
	if err != nil {
		return nil, err
	}

	spendables := sdk.Coins{}
	for _, denom := range lockedCoins.Denoms() {
		balance, err := bva.getBalance(ctx, accAddr, denom)
		if err != nil {
			return nil, err
		}

		lockedAmt := lockedCoins.AmountOf(balance.Denom)

		// get lockedCoin from that are not bonded for the sent denom
		notBondedLockedCoin, err := bva.GetNotBondedLockedCoin(ctx, sdk.NewCoin(balance.Denom, lockedAmt), balance.Denom)
		if err != nil {
			return nil, err
		}

		spendable, hasNeg := sdk.Coins{*balance}.SafeSub(notBondedLockedCoin)
		if hasNeg {
			spendable = sdk.Coins{sdk.NewCoin(balance.Denom, math.ZeroInt())}
		}

		spendables = spendables.Add(spendable...)
	}

	return &lockuptypes.QuerySpendableAmountResponse{
		SpendableTokens: spendables,
	}, nil
}

func (bva BaseLockup) RegisterExecuteHandlers(builder *accountstd.ExecuteBuilder) {
	accountstd.RegisterExecuteHandler(builder, bva.Undelegate)
	accountstd.RegisterExecuteHandler(builder, bva.WithdrawReward)
}

func (bva BaseLockup) RegisterQueryHandlers(builder *accountstd.QueryBuilder) {
	// accountstd.RegisterQueryHandler(builder, bva.QueryUnbondingEntries)
}
