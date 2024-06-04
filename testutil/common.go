package testutil

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/app"
)

// GeneratePrivKeyAddressPairsFromRand generates (deterministically) a total of n private keys and addresses.
func GeneratePrivKeyAddressPairs(n int) (keys []cryptotypes.PrivKey, addrs []sdk.AccAddress) {
	r := rand.New(rand.NewSource(12345)) // make the generation deterministic
	keys = make([]cryptotypes.PrivKey, n)
	addrs = make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		secret := make([]byte, 32)
		_, err := r.Read(secret)
		if err != nil {
			panic("Could not read randomness")
		}
		keys[i] = secp256k1.GenPrivKeyFromSecret(secret)
		addrs[i] = sdk.AccAddress(keys[i].PubKey().Address())
	}
	return
}

func InitSDKConfig() {
	// Set prefixes
	accountAddressPrefix := app.Bech32PrefixAccAddr
	accountPubKeyPrefix := app.Bech32PrefixAccPub
	validatorAddressPrefix := app.Bech32PrefixValAddr
	validatorPubKeyPrefix := app.Bech32PrefixValPub
	consNodeAddressPrefix := app.Bech32PrefixConsAddr
	consNodePubKeyPrefix := app.Bech32PrefixConsPub

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(accountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	// config.Seal()
}
