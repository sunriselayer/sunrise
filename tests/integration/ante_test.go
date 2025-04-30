package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/sunriselayer/sunrise/app/consts"
)

func TestAnteHandler(t *testing.T) {
	suite.Run(t, new(AnteHandlerTestSuite))
}

type AnteHandlerTestSuite struct {
	IntegrationTestSuite
}

func (s *AnteHandlerTestSuite) TestAnteHandlerWithValidFee() {
	// Create test accounts
	sender, senderPrivKey := s.createTestAccount("sender")
	recipient, _ := s.createTestAccount("recipient")

	// Build a transaction with valid fee
	msg := banktypes.NewMsgSend(sender, recipient, sdk.NewCoins(sdk.NewCoin(consts.FeeDenom, math.NewInt(100))))
	tx, err := s.buildTx(senderPrivKey, math.NewInt(1000), msg)
	require.NoError(s.T(), err)

	// Run the ante handler
	ctx := s.app.BaseApp.NewContext(false)
	anteHandler := s.app.AnteHandler()
	_, err = anteHandler(ctx, tx, false)
	require.NoError(s.T(), err)
}

func (s *AnteHandlerTestSuite) TestAnteHandlerWithInvalidFee() {
	// Create test accounts
	sender, senderPrivKey := s.createTestAccount("sender")
	recipient, _ := s.createTestAccount("recipient")

	// Build a transaction with invalid fee (too low)
	msg := banktypes.NewMsgSend(sender, recipient, sdk.NewCoins(sdk.NewCoin(consts.FeeDenom, math.NewInt(100))))
	tx, err := s.buildTx(senderPrivKey, math.NewInt(1), msg)
	require.NoError(s.T(), err)

	// Run the ante handler
	ctx := s.app.BaseApp.NewContext(false)
	anteHandler := s.app.AnteHandler()
	_, err = anteHandler(ctx, tx, false)
	require.Error(s.T(), err)
}

func (s *AnteHandlerTestSuite) TestAnteHandlerWithSwapBeforeFee() {
	// Create test accounts
	sender, senderPrivKey := s.createTestAccount("sender")
	recipient, _ := s.createTestAccount("recipient")

	// Build a transaction with SwapBeforeFee extension
	msg := banktypes.NewMsgSend(sender, recipient, sdk.NewCoins(sdk.NewCoin(consts.FeeDenom, math.NewInt(100))))
	tx, err := s.buildTx(senderPrivKey, math.NewInt(1000), msg)
	require.NoError(s.T(), err)

	// Run the ante handler
	ctx := s.app.BaseApp.NewContext(false)
	anteHandler := s.app.AnteHandler()
	_, err = anteHandler(ctx, tx, false)
	require.NoError(s.T(), err)
}
