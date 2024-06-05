BINARY=./build/sunrised
CHAIN_DIR=./data
CHAINID_1=test
NODE_HOME=$CHAIN_DIR/$CHAINID_1
BINARY_GOV_TOKEN=uvrise
BINARY_FEE_TOKEN=urise
VAL1=val
FAUCET=faucet
USER1=user1
USER2=user2
USER3=user3
USER4=user4

VAL_MNEMONIC_1="figure web rescue rice quantum sustain alert citizen woman cable wasp eyebrow monster teach hockey giant monitor hero oblige picnic ball never lamp distance"
FAUCET_MNEMONIC_1="chimney diesel tone pipe mouse detect vibrant video first jewel vacuum winter grant almost trim crystal similar giraffe dizzy hybrid trigger muffin awake leader"
USER_MNEMONIC_1="supply release type ostrich rib inflict increase bench wealth course enter pond spare neutral exact retire thing update inquiry atom health number lava taste"
USER_MNEMONIC_2="canyon second appear story film people resist slam waste again race rifle among room hip icon marriage sea quality prepare only liquid column click"
USER_MNEMONIC_3="among follow tooth egg unhappy city road expire solution visit visa skate allow network tissue slogan rose toddler develop utility negative peasant ostrich toward"
USER_MNEMONIC_4="charge split umbrella day gauge two orphan random human clerk buzz funny cabin purse fluid lecture blouse keen twist loud animal supply hat scare"

VAL_ADDRESS_1=sunrise1a8jcsmla6heu99ldtazc27dna4qcd4jyvln5d8
FAUCET_ADDRESS_1=sunrise1d6zd6awgjxuwrf4y863c9stz9m0eec4gnxuf79
USER_ADDRESS_1=sunrise155u042u8wk3al32h3vzxu989jj76k4zcc6d03n
USER_ADDRESS_2=sunrise1v0h8j7x7kfys29kj4uwdudcc9y0nx6tw2f955q
USER_ADDRESS_3=sunrise1y3t7sp0nfe2nfda7r9gf628g6ym6e7d43k5288
USER_ADDRESS_4=sunrise1pp2ruuhs0k7ayaxjupwj4k5qmgh0d72w8zu30p

conf="--home=$NODE_HOME --chain-id=$CHAINID_1 --keyring-backend=test -y --broadcast-mode=sync"
# conf="--home=$NODE_HOME --chain-id=$CHAINID_1 --keyring-backend=test -y --broadcast-mode=sync | grep txhash | awk '{ print $2 }'| xargs -I {} sh -c 'sleep 5; $0 q tx {}' $BINARY"

SSH_PREV_KEY_LOCATION=/your/ssh/key/location
