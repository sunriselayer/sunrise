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
# rm -rf $NODE_HOME &> /dev/null
$BINARY comet unsafe-reset-all --home $NODE_HOME --keep-addr-book

# Add directories for both chains, exit if an error occurs
if ! mkdir -p $NODE_HOME 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

rm -rf $NODE_HOME/config/genesis.json
cp $PROJECT_ROOT/build/genesis0731/genesis.json $NODE_HOME/config/genesis.json

# echo "Initializing $CHAINID_1..."
# $BINARY init test --home $NODE_HOME --chain-id=$CHAINID_1

# echo "Adding genesis accounts..."
# echo $VAL_MNEMONIC_1    | $BINARY keys add $VAL1 --home $NODE_HOME --recover --keyring-backend=test
# echo $FAUCET_MNEMONIC_1 | $BINARY keys add $FAUCET --home $NODE_HOME --recover --keyring-backend=test
# echo $USER_MNEMONIC_1 | $BINARY keys add $USER1 --home $NODE_HOME --recover --keyring-backend=test
# echo $USER_MNEMONIC_2 | $BINARY keys add $USER2 --home $NODE_HOME --recover --keyring-backend=test
# echo $USER_MNEMONIC_3 | $BINARY keys add $USER3 --home $NODE_HOME --recover --keyring-backend=test
# echo $USER_MNEMONIC_4 | $BINARY keys add $USER4 --home $NODE_HOME --recover --keyring-backend=test

# $BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $VAL1 --keyring-backend test -a) 100000000000000$BINARY_GOV_TOKEN,400000000000000$BINARY_NATIVE_TOKEN,100000000000$BINARY_STABLE_TOKEN --home $NODE_HOME

# echo "Creating and collecting gentx..."
# $BINARY genesis gentx $VAL1 7000000000$BINARY_GOV_TOKEN --home $NODE_HOME --chain-id $CHAINID_1 --keyring-backend test
# $BINARY genesis collect-gentxs --home $NODE_HOME

echo "Adding custom gentx files... from $PROJECT_ROOT/build/genesis0731/gentx"
GENTX_DIR="$PROJECT_ROOT/build/genesis0731/gentx"
JSON_FILES=("$GENTX_DIR"/*.json)
if [ -f "${JSON_FILES[0]}" ]; then
    echo "Found gentx files in $GENTX_DIR, adding to genesis.json..."
    jq -s '.[0].app_state.genutil.gen_txs += .[1:] | .[0]' "$NODE_HOME/config/genesis.json" "$GENTX_DIR"/*.json > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
fi

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

# jq ".app_state.gov.params.voting_period = \"60s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
# jq ".app_state.da.params.challenge_period = \"120s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
# jq ".app_state.da.params.rejected_removal_period = \"360s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
# jq ".app_state.da.params.verified_removal_period  = \"360s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
# jq ".app_state.da.params.proof_period  = \"120s\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;

# echo "Enable stable allowed addresses for user1"
# jq ".app_state.stable.params.allowed_addresses = [\"$($BINARY --home $NODE_HOME keys show $USER1 --keyring-backend test -a)\"]" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;

# echo "Enable tokenconverter allowed addresses"
# jq ".app_state.tokenconverter.params.allowed_addresses = [\"$($BINARY --home $NODE_HOME keys show $USER1 --keyring-backend test -a)\",\"$($BINARY --home $NODE_HOME keys show $FAUCET --keyring-backend test -a)\"]" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;

# echo "Enable fee burn"
# jq ".app_state.fee.params.burn_pool_id = \"1\"" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;
# jq ".app_state.fee.params.burn_enabled = true" $NODE_HOME/config/genesis.json > temp.json ; mv temp.json $NODE_HOME/config/genesis.json;

# echo "Enable urise send"
# jq '(.app_state.bank.send_enabled[] | select(.denom=="urise").enabled) = true' $NODE_HOME/config/genesis.json > $NODE_HOME/config/tmp_genesis.json && mv $NODE_HOME/config/tmp_genesis.json $NODE_HOME/config/genesis.json

echo "Add core-validator account"
jq --slurpfile coreValidatorAccounts "$PROJECT_ROOT/build/genesis0731/core-team-validator/accounts_core_validator.json" '.app_state.auth.accounts += $coreValidatorAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile coreValidatorBalance "$PROJECT_ROOT/build/genesis0731/core-team-validator/balances_core_validator.json" '.app_state.bank.balances += $coreValidatorBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add validator accounts"
jq --slurpfile validatorAccounts "$PROJECT_ROOT/build/genesis0731/validators/accounts_validator.json" '.app_state.auth.accounts += $validatorAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile validatorBalance "$PROJECT_ROOT/build/genesis0731/validators/balances_validator.json" '.app_state.bank.balances += $validatorBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/genesis0731/validators/init_lockup_msgs_validator.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add airdrop accounts"
jq --slurpfile airdropAccounts "$PROJECT_ROOT/build/genesis0731/airdrop/accounts_airdrop_exclude_4vals.json" '.app_state.auth.accounts += $airdropAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile airdropBalance "$PROJECT_ROOT/build/genesis0731/airdrop/balances_airdrop_exclude_4vals_add_public_sale.json" '.app_state.bank.balances += $airdropBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/genesis0731/airdrop/init_lockup_msgs_airdrop_include_4vals.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add partner accounts"
jq --slurpfile partnerAccounts "$PROJECT_ROOT/build/genesis0731/partners/accounts_partner.json" '.app_state.auth.accounts += $partnerAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile partnerBalance "$PROJECT_ROOT/build/genesis0731/partners/balances_partner.json" '.app_state.bank.balances += $partnerBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/genesis0731/partners/init_lockup_msgs_partner.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add faucet accounts"
jq --slurpfile faucetAccounts "$PROJECT_ROOT/build/genesis0731/faucets/accounts_faucet.json" '.app_state.auth.accounts += $faucetAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile faucetBalance "$PROJECT_ROOT/build/genesis0731/faucets/balances_faucet.json" '.app_state.bank.balances += $faucetBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add public sale accounts"
jq --slurpfile publicSaleAccounts "$PROJECT_ROOT/build/genesis0731/public-sale/no-airdrop/accounts_public_sale.json" '.app_state.auth.accounts += $publicSaleAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile publicSaleBalance "$PROJECT_ROOT/build/genesis0731/public-sale/no-airdrop/balances_public_sale.json" '.app_state.bank.balances += $publicSaleBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/genesis0731/public-sale/no-airdrop/init_lockup_msgs_public_sale_no_airdrop.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/genesis0731/public-sale/no-airdrop/init_lockup_msgs_public_sale_no_airdrop_bonus.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/genesis0731/public-sale/airdrop/init_lockup_msgs_public_sale_airdrop.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/genesis0731/public-sale/airdrop/init_lockup_msgs_public_sale_airdrop_bonus.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add ununifi validator accounts"
jq --slurpfile ununifiValidatorAccounts "$PROJECT_ROOT/build/genesis0731/ununifi-validators/accounts_ununifi_validator.json" '.app_state.auth.accounts += $ununifiValidatorAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile ununifiValidatorBalance "$PROJECT_ROOT/build/genesis0731/ununifi-validators/balances_ununifi_validator.json" '.app_state.bank.balances += $ununifiValidatorBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/genesis0731/ununifi-validators/init_lockup_msgs_ununifi_validator.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add advisor accounts"
jq --slurpfile advisorAccounts "$PROJECT_ROOT/build/genesis0731/advisors/accounts_advisor.json" '.app_state.auth.accounts += $advisorAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile advisorBalance "$PROJECT_ROOT/build/genesis0731/advisors/balances_advisor.json" '.app_state.bank.balances += $advisorBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/genesis0731/advisors/init_lockup_msgs_advisor.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add team accounts"
jq --slurpfile teamAccounts "$PROJECT_ROOT/build/genesis0731/team/accounts_team.json" '.app_state.auth.accounts += $teamAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile teamBalance "$PROJECT_ROOT/build/genesis0731/team/balances_team.json" '.app_state.bank.balances += $teamBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile lockupMsgs "$PROJECT_ROOT/build/genesis0731/team/init_lockup_msgs_team.json" '.app_state.lockup.init_lockup_msgs += $lockupMsgs[0].init_lockup_msgs' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"

echo "Add foundation accounts"
jq --slurpfile foundationAccounts "$PROJECT_ROOT/build/genesis0731/foundation/accounts_foundation.json" '.app_state.auth.accounts += $foundationAccounts[0].accounts' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
jq --slurpfile foundationBalance "$PROJECT_ROOT/build/genesis0731/foundation/balances_foundation.json" '.app_state.bank.balances += $foundationBalance[0].balances' "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"


echo "Update supply"
jq ".app_state.bank.supply = [{\"denom\": \"urise\", \"amount\": \"449999965000000\"}, {\"denom\": \"uusdrise\", \"amount\": \"10003500000\"}, {\"denom\": \"uvrise\", \"amount\": \"50000035000000\"}]" "$NODE_HOME/config/genesis.json" > temp.json && mv temp.json "$NODE_HOME/config/genesis.json"
# Register accounts after genesis
# $BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $FAUCET --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000000000$BINARY_NATIVE_TOKEN,100000000000$BINARY_STABLE_TOKEN,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
# $BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER1 --keyring-backend test -a) 100000000000$BINARY_NATIVE_TOKEN --home $NODE_HOME
# $BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER2 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_NATIVE_TOKEN,100000000000$BINARY_STABLE_TOKEN,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
# $BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER3 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_NATIVE_TOKEN,100000000000$BINARY_STABLE_TOKEN,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
# $BINARY genesis add-genesis-account $($BINARY --home $NODE_HOME keys show $USER4 --keyring-backend test -a) 100000000000$BINARY_GOV_TOKEN,100000000000$BINARY_NATIVE_TOKEN,100000000000$BINARY_STABLE_TOKEN,100000000000000uusdt,100000000000000uusdc --home $NODE_HOME
