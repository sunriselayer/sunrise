package types

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidatorSeed(t *testing.T) {
	valAddr1 := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
	valAddr2 := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
	valAddr3 := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
	valAddr4 := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
	valAddr5 := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())

	assert.NotPanics(t, func() { ValidatorSeed(valAddr1) }, "Assert not panic")
	assert.NotPanics(t, func() { ValidatorSeed(valAddr2) }, "Assert not panic")
	assert.NotPanics(t, func() { ValidatorSeed(valAddr3) }, "Assert not panic")
	assert.NotPanics(t, func() { ValidatorSeed(valAddr4) }, "Assert not panic")
	assert.NotPanics(t, func() { ValidatorSeed(valAddr5) }, "Assert not panic")
}
