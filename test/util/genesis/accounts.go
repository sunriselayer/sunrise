package genesis

import (
	"context"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/sunriselayer/sunrise-app/app"
	"github.com/sunriselayer/sunrise-app/app/encoding"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Account struct {
	Name          string
	InitialTokens int64
}

func NewAccounts(initBal int64, names ...string) []Account {
	accounts := make([]Account, len(names))
	for i, name := range names {
		accounts[i] = Account{
			Name:          name,
			InitialTokens: initBal,
		}
	}
	return accounts
}

func (ga *Account) ValidateBasic() error {
	if ga.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if ga.InitialTokens <= 0 {
		return fmt.Errorf("initial tokens must be positive")
	}
	return nil
}

type Validator struct {
	Account
	Stake int64

	// ConsensusKey is the key used by the validator to sign votes.
	ConsensusKey crypto.PrivKey
	NetworkKey   crypto.PrivKey
}

func NewDefaultValidator(name string) Validator {
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	return Validator{
		Account: Account{
			Name:          name,
			InitialTokens: 999_999_999_999_999_999,
		},
		Stake:        99_999_999_999_999_999, // save some tokens for fees
		ConsensusKey: GenerateEd25519(NewSeed(r)),
		NetworkKey:   GenerateEd25519(NewSeed(r)),
	}
}

// ValidateBasic performs stateless validation on the validitor
func (v *Validator) ValidateBasic() error {
	if err := v.Account.ValidateBasic(); err != nil {
		return err
	}
	if v.Stake <= 0 {
		return fmt.Errorf("stake must be positive")
	}
	if v.ConsensusKey == nil {
		return fmt.Errorf("consensus key cannot be empty")
	}
	if v.Stake > v.InitialTokens {
		return fmt.Errorf("stake cannot be greater than initial tokens")
	}
	return nil
}

// GenTx generates a genesis transaction to create a validator as configured by
// the validator struct. It assumes the validator's genesis account has already
// been added to the keyring and that the sequence for that account is 0.
func (v *Validator) GenTx(ctx context.Context, ecfg encoding.Config, kr keyring.Keyring, chainID string) (sdk.Tx, error) {
	rec, err := kr.Key(v.Name)
	if err != nil {
		return nil, err
	}
	addr, err := rec.GetAddress()
	if err != nil {
		return nil, err
	}

	commission, err := sdkmath.LegacyNewDecFromStr("0.5")
	if err != nil {
		return nil, err
	}

	pk, err := cryptocodec.FromTmPubKeyInterface(v.ConsensusKey.PubKey())
	if err != nil {
		return nil, fmt.Errorf("converting public key for node %s: %w", v.Name, err)
	}

	createValMsg, err := stakingtypes.NewMsgCreateValidator(
		sdk.ValAddress(addr).String(),
		pk,
		sdk.NewCoin(app.BondDenom, sdkmath.NewInt(v.Stake)),
		stakingtypes.NewDescription(v.Name, "", "", "", ""),
		stakingtypes.NewCommissionRates(commission, sdkmath.LegacyOneDec(), sdkmath.LegacyOneDec()),
		sdkmath.NewInt(v.Stake/2),
	)
	if err != nil {
		return nil, err
	}

	fee := sdk.NewCoins(sdk.NewCoin(app.BondDenom, sdkmath.NewInt(1)))
	txBuilder := ecfg.TxConfig.NewTxBuilder()
	err = txBuilder.SetMsgs(createValMsg)
	if err != nil {
		return nil, err
	}
	txBuilder.SetFeeAmount(fee)    // Arbitrary fee
	txBuilder.SetGasLimit(1000000) // Need at least 100386

	txFactory := tx.Factory{}
	txFactory = txFactory.
		WithChainID(chainID).
		WithKeybase(kr).
		WithTxConfig(ecfg.TxConfig)

	err = tx.Sign(ctx, txFactory, v.Name, txBuilder, true)
	if err != nil {
		return nil, err
	}

	return txBuilder.GetTx(), nil
}
