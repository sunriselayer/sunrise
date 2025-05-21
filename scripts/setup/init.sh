#!/bin/bash

# Load shell variables
SCRIPT_DIR=$(cd $(dirname $0); pwd)
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

$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $VAL1 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_FEE_TOKEN --home $NODE_HOME
$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $FAUCET --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_FEE_TOKEN,100000000000000uglu,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER1 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_FEE_TOKEN,100000000000000uglu,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER2 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_FEE_TOKEN,100000000000000uglu,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER3 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_FEE_TOKEN,100000000000000uglu,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
$BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER4 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_FEE_TOKEN,100000000000000uglu,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME

echo "Creating and collecting gentx..."
$BINARY genesis gentx $VAL1 7000000000$BINARY_GOV_TOKEN --home $NODE_HOME --chain-id $CHAINID_1 --keyring-backend test --fees 100000$BINARY_FEE_TOKEN
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

jq ".app_state.gov.params.voting_period = \"20s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
jq ".consensus.params.feature.vote_extensions_enable_height = \"1\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
jq ".app_state.da.params.challenge_period = \"30s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
jq ".app_state.da.params.rejected_removal_period = \"60s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
jq ".app_state.da.params.verified_removal_period  = \"60s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
# jq ".app_state.accounts.init_account_msgs  = [{\"sender\":\"sunrise1d6zd6awgjxuwrf4y863c9stz9m0eec4gnxuf79\",\"account_type\":\"self-delegatable-continuous-locking-account\",\"message\": {\"@type\": \"/sunrise.accounts.self_delegatable_lockup.v1.MsgInitSelfDelegatableLockupAccount\",\"owner\": \"sunrise155u042u8wk3al32h3vzxu989jj76k4zcc6d03n\",\"end_time\": \"2026-01-01T00:00:00Z\",\"start_time\": \"2025-01-01T00:00:00Z\"}, \"funds\": [{\"denom\": \"urise\", \"amount\": \"5000000\"}]}]" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
