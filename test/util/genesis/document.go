package genesis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sunriselayer/sunrise/app/consts"
	"github.com/sunriselayer/sunrise/app/encoding"

	sdkmath "cosmossdk.io/math"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	coretypes "github.com/cometbft/cometbft/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	testencoding "github.com/sunriselayer/sunrise/test/util/encoding"
)

// Document will create a valid genesis doc with funded addresses.
func Document(
	ecfg encoding.Config,
	params *cmtproto.ConsensusParams,
	chainID string,
	gentxs []json.RawMessage,
	addrs []string,
	pubkeys []cryptotypes.PubKey,
	mods ...Modifier,
) (*coretypes.GenesisDoc, error) {
	genutilGenState := genutiltypes.DefaultGenesisState()
	genutilGenState.GenTxs = gentxs

	genBals, genAccs, err := accountsToSDKTypes(addrs, pubkeys)
	if err != nil {
		return nil, err
	}

	accounts, err := authtypes.PackAccounts(genAccs)
	if err != nil {
		return nil, err
	}

	authGenState := authtypes.DefaultGenesisState()
	bankGenState := banktypes.DefaultGenesisState()
	authGenState.Accounts = append(authGenState.Accounts, accounts...)
	bankGenState.Balances = append(bankGenState.Balances, genBals...)
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

	// perform some basic validation of the genesis state
	if err := authtypes.ValidateGenesis(*authGenState); err != nil {
		return nil, err
	}
	if err := bankGenState.Validate(); err != nil {
		return nil, err
	}
	// if err := genutiltypes.ValidateGenesis(genutilGenState, ecfg.TxConfig.TxJSONDecoder()); err != nil {
	// 	return nil, err
	// }

	state := testencoding.ModuleBasics.DefaultGenesis(ecfg.Codec)
	state[authtypes.ModuleName] = ecfg.Codec.MustMarshalJSON(authGenState)
	state[banktypes.ModuleName] = ecfg.Codec.MustMarshalJSON(bankGenState)
	state[genutiltypes.ModuleName] = ecfg.Codec.MustMarshalJSON(genutilGenState)

	for _, modifer := range mods {
		state = modifer(state)
	}

	stateBz, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return nil, err
	}

	// Create the genesis doc
	genesisDoc := &coretypes.GenesisDoc{
		ChainID:     chainID,
		GenesisTime: time.Now(),
		ConsensusParams: &coretypes.ConsensusParams{
			Block:     coretypes.BlockParams(*params.Block),
			Evidence:  coretypes.EvidenceParams(*params.Evidence),
			Validator: coretypes.ValidatorParams(*params.Validator),
			Version:   coretypes.VersionParams(*params.Version),
			ABCI:      coretypes.ABCIParams(*params.Abci),
		},
		AppState: stateBz,
	}

	return genesisDoc, nil
}

// accountsToSDKTypes converts the genesis accounts to native SDK types.
func accountsToSDKTypes(addrs []string, pubkeys []cryptotypes.PubKey) ([]banktypes.Balance, []authtypes.GenesisAccount, error) {
	if len(addrs) != len(pubkeys) {
		return nil, nil, fmt.Errorf("length of addresses and public keys are not equal")
	}
	genBals := make([]banktypes.Balance, len(addrs))
	genAccs := make([]authtypes.GenesisAccount, len(addrs))
	hasMap := make(map[string]bool)
	for i, addr := range addrs {
		if hasMap[addr] {
			return nil, nil, fmt.Errorf("duplicate account address %s", addr)
		}
		hasMap[addr] = true

		pubKey := pubkeys[i]

		balances := sdk.NewCoins(
			sdk.NewCoin(consts.BondDenom, sdkmath.NewInt(999_999_999_999_999_999)),
		)

		genBals[i] = banktypes.Balance{Address: addr, Coins: balances.Sort()}

		parsedAddress, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return nil, nil, err
		}

		genAccs[i] = authtypes.NewBaseAccount(parsedAddress, pubKey, uint64(i), 0)
	}
	return genBals, genAccs, nil
}
