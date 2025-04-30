package integration

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"

	"github.com/sunriselayer/sunrise/app"
	"github.com/sunriselayer/sunrise/app/consts"
)

// testAppOptions implements servertypes.AppOptions
type testAppOptions struct {
	minGasPrices string
}

func (ao testAppOptions) Get(key string) interface{} {
	switch key {
	case server.FlagMinGasPrices:
		return ao.minGasPrices
	default:
		return nil
	}
}

// IntegrationTestSuite represents the e2e test suite
type IntegrationTestSuite struct {
	suite.Suite
	ctx       context.Context
	cancel    context.CancelFunc
	app       *app.App
	kr        keyring.Keyring
	clientCtx client.Context
}

// SetupSuite runs once before all tests
func (s *IntegrationTestSuite) SetupSuite() {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// Initialize app
	appOptions := testAppOptions{
		minGasPrices: fmt.Sprintf("%v%s", consts.DefaultMinGasPrice, consts.FeeDenom),
	}

	s.app = app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		appOptions,
	)

	// Initialize keyring
	s.kr = keyring.NewInMemory(s.app.AppCodec())

	// Initialize client context
	s.clientCtx = client.Context{}.
		WithCodec(s.app.AppCodec()).
		WithInterfaceRegistry(s.app.InterfaceRegistry()).
		WithTxConfig(s.app.TxConfig()).
		WithKeyring(s.kr)
}

// createTestAccount creates a new test account with the given name and returns its address and private key
func (s *IntegrationTestSuite) createTestAccount(name string) (sdk.AccAddress, *secp256k1.PrivKey) {
	// Create a new private key
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()

	// Create a new account in the keyring
	record, err := s.kr.SaveOfflineKey(name, pubKey)
	s.Require().NoError(err)

	// Get the account address
	addr, err := record.GetAddress()
	s.Require().NoError(err)

	return addr, privKey
}

// buildTx builds a transaction with the given messages and fee amount
func (s *IntegrationTestSuite) buildTx(privKey *secp256k1.PrivKey, feeAmount math.Int, msgs ...sdk.Msg) (authsigning.Tx, error) {
	// Create a new transaction builder
	txBuilder := s.clientCtx.TxConfig.NewTxBuilder()

	// Set the messages
	err := txBuilder.SetMsgs(msgs...)
	s.Require().NoError(err)

	// Set the gas limit
	txBuilder.SetGasLimit(200000)

	// Set the fee
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(consts.FeeDenom, feeAmount)))

	// Set the memo
	txBuilder.SetMemo("test transaction")

	// Set the timeout height
	txBuilder.SetTimeoutHeight(0)

	// Create the signer data
	signerData := authsigning.SignerData{
		ChainID:       "test-chain",
		AccountNumber: 0,
		Sequence:      0,
	}

	// Sign the transaction
	sig, err := tx.SignWithPrivKey(
		s.ctx,
		signing.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder,
		privKey,
		s.clientCtx.TxConfig,
		0,
	)
	s.Require().NoError(err)

	// Set the signature
	err = txBuilder.SetSignatures(sig)
	s.Require().NoError(err)

	return txBuilder.GetTx(), nil
}
