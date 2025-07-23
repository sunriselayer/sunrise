#!/bin/bash

# Load shell variables
SCRIPT_DIR=$(cd $(dirname $0); pwd)
PROJECT_ROOT=$(cd "$SCRIPT_DIR/../.."; pwd)
. $SCRIPT_DIR/variables.sh

# Stop if it is already running 
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    killall $BINARY
fi

echo "Removing previous data..."
rm -rf $NODE_HOME &> /dev/null

# Add directories for both chains, exit if an error occurs
if ! mkdir -p $NODE_HOME 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

echo "Initializing $CHAINID_1..."
$BINARY init test --home $NODE_HOME --chain-id=$CHAINID_1

echo "Adding genesis accounts..."
echo $VAL_MNEMONIC_1    | $BINARY keys add $VAL1 --home $NODE_HOME --recover --keyring-backend=test
echo $FAUCET_MNEMONIC_1 | $BINARY keys add $FAUCET --home $NODE_HOME --recover --keyring-backend=test
echo $USER_MNEMONIC_1 | $BINARY keys add $USER1 --home $NODE_HOME --recover --keyring-backend=test
echo $USER_MNEMONIC_2 | $BINARY keys add $USER2 --home $NODE_HOME --recover --keyring-backend=test
echo $USER_MNEMONIC_3 | $BINARY keys add $USER3 --home $NODE_HOME --recover --keyring-backend=test
echo $USER_MNEMONIC_4 | $BINARY keys add $USER4 --home $NODE_HOME --recover --keyring-backend=test

$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $VAL1 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_NATIVE_TOKEN,100000000000$BINARY_STABLE_TOKEN --home $NODE_HOME
$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $FAUCET --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000000000$BINARY_NATIVE_TOKEN,100000000000$BINARY_STABLE_TOKEN,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER1 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_NATIVE_TOKEN,100000000000$BINARY_STABLE_TOKEN,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER2 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_NATIVE_TOKEN,100000000000$BINARY_STABLE_TOKEN,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER3 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_NATIVE_TOKEN,100000000000$BINARY_STABLE_TOKEN,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER4 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_NATIVE_TOKEN,100000000000$BINARY_STABLE_TOKEN,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME

echo "Creating and collecting gentx..."
$BINARY genesis gentx $VAL1 7000000000$BINARY_GOV_TOKEN --home $NODE_HOME --chain-id $CHAINID_1 --keyring-backend test
$BINARY genesis collect-gentxs --home $NODE_HOME

echo "Changing defaults config files..."
OS=$(uname -s)
if [ "$OS" == "Darwin" ]; then
  echo $OS
  sleep 1
  sed_i="sed -i '' "
elif [ "$OS" == "Linux" ]; then
  echo $OS
  sleep 1
  sed_i="sed -i"
fi

$sed_i '/\[api\]/,+3 s/enable = false/enable = true/' $NODE_HOME/config/app.toml;
$sed_i 's/mode = "full"/mode = "validator"/' $NODE_HOME/config/config.toml;
$sed_i "s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/" $NODE_HOME/config/app.toml;
$sed_i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/' $NODE_HOME/config/config.toml;

jq ".app_state.gov.params.voting_period = \"60s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
jq ".app_state.da.params.challenge_period = \"120s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
jq ".app_state.da.params.rejected_removal_period = \"360s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
jq ".app_state.da.params.verified_removal_period  = \"360s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
jq ".app_state.da.params.proof_period  = \"120s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;

echo "Enable stable allowed addresses for user1"
jq ".app_state.stable.params.allowed_addresses = [\"$($BINARY --home $NODE_HOME keys show $USER1 --keyring-backend test -a)\"]" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;

echo "Enable tokenconverter allowed addresses for user1"
jq ".app_state.tokenconverter.params.allowed_addresses = [\"$($BINARY --home $NODE_HOME keys show $USER1 --keyring-backend test -a)\"]" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;

echo "Enable fee burn"
jq ".app_state.fee.params.burn_pool_id = \"1\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
jq ".app_state.fee.params.burn_enabled = true" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;

# echo "Enable urise send"
# jq '(.app_state.bank.send_enabled[] | select(.denom=="urise").enabled) = true' $NODE_HOME/config/genesis.json > $NODE_HOME/config/tmp_genesis.json && mv $NODE_HOME/config/tmp_genesis.json $NODE_HOME/config/genesis.json


echo "Add airdrop accounts"
jq --slurpfile airdropAccounts "$PROJECT_ROOT/build/genesis_accounts_airdrop.json" '.app_state.auth.accounts += $airdropAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile airdropBalance "$PROJECT_ROOT/build/genesis_balances_airdrop_v2.json" '.app_state.bank.balances += $airdropBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add init lockup messages for airdrop from json file"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/init_lockup_msgs_airdrop.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add validator accounts"
jq --slurpfile validatorAccounts "$PROJECT_ROOT/build/genesis_accounts_validator.json" '.app_state.auth.accounts += $validatorAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile validatorBalance "$PROJECT_ROOT/build/genesis_balances_validator.json" '.app_state.bank.balances += $validatorBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add partner accounts"
jq --slurpfile partnerAccounts "$PROJECT_ROOT/build/genesis_accounts_partner.json" '.app_state.auth.accounts += $partnerAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile partnerBalance "$PROJECT_ROOT/build/genesis_balances_partner.json" '.app_state.bank.balances += $partnerBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add init lockup messages for partner from json file"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/init_lockup_msgs_partner.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Update balances for a specific address"
#jq '(.app_state.bank.balances[] | select(.address == "sunrise1d6zd6awgjxuwrf4y863c9stz9m0eec4gnxuf79").coins[] | select(.denom == "urise")).amount = "99995635386493883"' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq '(.app_state.bank.balances[] | select(.address == "sunrise1d6zd6awgjxuwrf4y863c9stz9m0eec4gnxuf79").coins[] | select(.denom == "urise")).amount = "99991887886493883"' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq '(.app_state.bank.balances[] | select(.address == "sunrise1d6zd6awgjxuwrf4y863c9stz9m0eec4gnxuf79").coins[] | select(.denom == "uusdrise")).amount = "99996700000"' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq '(.app_state.bank.balances[] | select(.address == "sunrise1d6zd6awgjxuwrf4y863c9stz9m0eec4gnxuf79").coins[] | select(.denom == "uvrise")).amount = "99967000000"' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
