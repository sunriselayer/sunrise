#!/bin/sh



rm -rf $HOME/.sunrise/



chain_id=sunrise-1



sunrised init --chain-id=$chain_id sunrise-test --home=$HOME/.sunrise

sunrised keys add validator --keyring-backend=test --home=$HOME/.sunrise



echo "stuff floor burst decline forum faith shoe crack sponsor palm blossom dirt awake protect gorilla fashion insane soup raise coffee work mixture bracket hunt" > validator2.txt;

sunrised keys add validator2 --keyring-backend=test --home=$HOME/.sunrise --recover < validator2.txt;



VALIDATOR=$(sunrised keys show validator -a --keyring-backend=test --home=$HOME/.sunrise)

sunrised genesis add-genesis-account $VALIDATOR 1000000000000000000000000000000tttt,100000000000000000000000000000uuuu,100000000000uvrise --home=$HOME/.sunrise

sunrised genesis gentx validator 500000000uvrise --keyring-backend=test --home=$HOME/.sunrise --chain-id=$chain_id

sunrised genesis collect-gentxs --home=$HOME/.sunrise



sed -i '' 's#"vote_extensions_enable_height": "0"#"vote_extensions_enable_height": "1"#g' $HOME/.sunrise/config/genesis.json

sed -i '' 's#"voting_period": "604800s"#"voting_period": "30s"#g' $HOME/.sunrise/config/genesis.json

sed -i '' '160s/false/true/' $HOME/.sunrise/config/app.toml

sunrised start --home=$HOME/.sunrise --minimum-gas-prices=0uvrise --log_level=debug



# sunrised tx bank send validator sunrise10wl2y0h3cm6533u3k0sg2eqqd8s4npp05mqz5r 1000000uuuu --keyring-backend=test --chain-id=sunrise-1 --yes --broadcast-mode=sync

# sunrised query bank balances sunrise1fh7903espudmcguqvd2qpn73ym5wftg8jd685v



sunrised keys show -a --keyring-backend=test validator

sunrised query bank balances sunrise10wl2y0h3cm6533u3k0sg2eqqd8s4npp05mqz5r



sunrised tx liquiditypool create-pool --denom-base=tttt --denom-quote=uuuu --fee-rate="0.01" --price-ratio="1.000100000000000000" --base-offset="0.500000000000000000" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-pool --denom-base=tttt --denom-quote=uuuu --fee-rate="0.01" --price-ratio="1.000100000000000000" --base-offset="0.000000000000000000" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"



sunrised query liquiditypool list-pools

# Pool:

# - current_sqrt_price: "0"

#   current_tick_liquidity: "0"

#   denom_base: tttt

#   denom_quote: uuuu

#   fee_rate: "0"

#   tick_params:

#     base_offset: "500000000000000000"

#     price_ratio: "1000100000000000000"



sunrised tx liquiditypool create-position --pool-id=1 --lower-tick="-10" --upper-tick=10 --token-base=10000000tttt --token-quote=10000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-position --pool-id=1 --lower-tick="-20" --upper-tick=20 --token-base=10000000tttt --token-quote=10000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-position --pool-id=1 --lower-tick="-20" --upper-tick=10 --token-base=10000000tttt --token-quote=10000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-position --pool-id=1 --lower-tick="100" --upper-tick=200 --token-base=30000000tttt --token-quote=20000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-position --pool-id=1 --lower-tick="-200" --upper-tick="-100" --token-base=30000000tttt --token-quote=20000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-position --pool-id=1 --lower-tick="-830000" --upper-tick=830000 --token-base=10000000tttt --token-quote=10000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-position --pool-id=0 --lower-tick="-830000" --upper-tick=830000 --token-base=100000000000000000000000tttt --token-quote=100uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"



# out of bound ticks

sunrised tx liquiditypool create-position --pool-id=1 --lower-tick="-500000" --upper-tick=800000 --token-base=3000000000tttt --token-quote=1000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-position --pool-id=1 --lower-tick="-400000" --upper-tick=1000000 --token-base=3000000000tttt --token-quote=1000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"



sunrised tx liquiditypool create-position --pool-id=2 --lower-tick="-10" --upper-tick=10 --token-base=10000000tttt --token-quote=10000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-position --pool-id=2 --lower-tick="-20" --upper-tick=20 --token-base=10000000tttt --token-quote=10000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-position --pool-id=2 --lower-tick="-20" --upper-tick=10 --token-base=10000000tttt --token-quote=10000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-position --pool-id=2 --lower-tick="100" --upper-tick=200 --token-base=30000000tttt --token-quote=20000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool create-position --pool-id=2 --lower-tick="-200" --upper-tick="-100" --token-base=30000000tttt --token-quote=20000000uuuu --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"



sunrised query liquiditypool list-positions

sunrised tx liquiditypool increase-liquidity --id="1" --amount-base=10000000 --amount-quote=10000000 --min-amount-base="0" --min-amount-quote="0" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --gas=1000000 --fees="10000uvrise"

sunrised tx liquiditypool decrease-liquidity --id="2" --liquidity=10005500824958739014503444368 --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool decrease-liquidity --id="0" --liquidity=14772255750516584004256859 --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"





sunrised tx swap swap-exact-amount-in --interface-provider="sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f" --route='{"denom_in":"tttt","denom_out":"uuuu","pool":{"pool_id":1}}' --amount-in="10000000" --min-amount-out=1000 --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx swap swap-exact-amount-out --interface-provider="sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f" --route='{"denom_in":"uuuu","denom_out":"tttt","pool":{"pool_id":1}}' --max-amount-in="1000000000" --amount-out=5000000 --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx swap swap-exact-amount-out --interface-provider="sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f" --route='{"denom_in":"tttt","denom_out":"uuuu","pool":{"pool_id":1}}' --max-amount-in="1000000000000000000000000" --amount-out=1 --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"



sunrised tx swap swap-exact-amount-in --interface-provider="sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f" --route='{"denom_in":"tttt","denom_out":"uuuu","pool":{"pool_id":2}}' --amount-in="10000000" --min-amount-out=1000 --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx swap swap-exact-amount-out --interface-provider="sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f" --route='{"denom_in":"uuuu","denom_out":"tttt","pool":{"pool_id":2}}' --max-amount-in="100000000" --amount-out=10000000 --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"



sunrised query liquiditypool show-position-fees 0

sunrised query liquiditypool show-position-fees 1

sunrised query liquiditypool show-position-fees 2

sunrised query liquiditypool show-position-fees 3

sunrised query liquiditypool show-position-fees 4

sunrised query liquiditypool show-position-fees 5

sunrised query liquiditypool show-position-fees 6

sunrised query liquiditypool show-position-fees 7

sunrised query liquiditypool show-position-fees 8

sunrised query liquiditypool show-position-fees 9

sunrised query liquiditypool show-position-fees 10

sunrised tx liquiditypool claim-rewards --position-ids=0 --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool claim-rewards --position-ids=0,1 --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquiditypool claim-rewards --position-ids=1,2 --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"



sunrised tx liquidityincentive vote-gauge --weights='{"pool_id":0,"weight":"1000000000000000000"}' --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"

sunrised tx liquidityincentive vote-gauge --weights='{"pool_id":0,"weight":"500000000000000000"}' --weights='{"pool_id":1,"weight":"500000000000000000"}' --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"



sunrised query liquidityincentive list-epoch

sunrised query liquidityincentive list-gauge 0

# sunrised export  --home=$HOME/.sunrise



sunrised tx da publish-data --recovered-data-hash="010101" --metadata-uri="metadata" --shard-double-hashes="010101" --from=validator --chain-id=sunrise-1 --yes --keyring-backend=test --fees="10000uvrise"



sunrised query da published-data