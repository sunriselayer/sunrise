package v0_2_0_test

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func upgradeSendCoin(
	ctx sdk.Context,
	bankkeeper bankkeeper.Keeper,
) error {
	fromAddress := "sunrise1kw8x5dncdw7ualrx02q4cldcxhsmg5vwtxaxvq" // core dev-2
	toAddresses := []string{
		// new validators
		"sunrise1t43ckph4y0vg4854ktms29nr6ta6dnlprlm756", // Neukind
		"sunrise17pgafjy76nwda34ffxqwrc32vqutapdkhtlqf5", // Stakelab
		"sunrise1cp3wxy6d3kjlt850rh9pjf0ma804fghpjx4gpn", // TC Network
		"sunrise1lz4rtdlgvql0y07g9tyvvuhpxt42pd5qe8ee7d", // Krewz
		"sunrise1gx63mdytxuej9s3lk3n9hkw6mhuhvphrq4tcp8", // Dwentz
		"sunrise1ztckrcazwjf40uskh6zyhnvgrulykd67yeetlc", // Rhino
		"sunrise1dzuvf3dd83csf9d5wyqwl5evkxe6fwk48m0u5g", // deNodes
		"sunrise13f85s2dgane2kmx9e032s47w7f0xk3rz38x49x", // Kocality
		"sunrise1xzern65ppact6sgh4htllnr8pg9dpnc836jt62", // SynergyNodes
		"sunrise1fgpa82a0fzn8vkj37uzwz54u20peevh4nt6e5z", // Autostake
		"sunrise1u8hvyz3wqe3rm4zk9zw6upsn2f3nk8k3r6vzrh", // Decentrio
		"sunrise12y4nkeug9j699jcrg28v2u6hx88l29h3vqgner", // Bware Labs
		"sunrise12q0e8sya2j0lmzcya4jhecyjl47v0xuj256gg2", // STC
		"sunrise14qekdkj2nmmwea4ufg9n002a3pud23y8h0t6qr", // Lavender.five
		"sunrise1a5t392uyw8x0dmul48lfrt6n7emvmzt0sdlllu", // POSTHUMAN
		"sunrise1j4837fmd8zxdeypv8ejk6mkwav5jtvglflny3d", // Staking4All
		"sunrise13xhjcwatxvjpdsx5lpk7t9wryhxr7wgj4ywu6n", // Coinhunters
		// "sunrise1etx55kw7tkmnjqz0k0mups4ewxlr324t8c7704", // Nodestake (already)
		"sunrise1u2f59za8eh2lj4n42rx9yyjcdrwwufquy0fenx", // Anode.team
		"sunrise16kjq7nd2235fnamruptm59ma5y69p90rcrfrwv", // DAIC
		"sunrise1mwrm63jthf0frsqkx9em3mzkw2fsu3nqsc02um", // Nodes.Guru
		"sunrise1knkwlqqxky07f7vu6vgxna49m2tf9aa9jlwtrw", // [NODERS]
		"sunrise1mjkhjz9y0mhnfwyma9nkh4ufr4n4ms2gsy5lxz", // Stake-take
		"sunrise16kjq7nd2235fnamruptm59ma5y69p90rcrfrwv", // Coinage x DAIC
		// "sunrise1yvh0yvum6t2rc9xsv9cdw26xn393q5a0wxreq3", // HashKey Cloud (already)
		"sunrise17prvyymh0tf25aee5fazjxtdfs5vex2w5fsk3f", // Citadel.one
		// "sunrise133fhqjfd3mmku3vwczsf5lyzuhavfmqly4f79w", // HoodRun (already)
		"sunrise1ezwcuz7qc5t68650j3cj4uwpxhkm05pztdu83f", // Polkachu
		"sunrise13htn2mkup5kpkaw00ss9lzscpf4k3ppwk0z5hz", // MeKongLabs
		"sunrise1qxvq5rc43umg96hr2enup5ywd20lzgyhtd3hts", // KysenPool
		"sunrise1y6phqcxhaka6ax4trr8wjrepp97ag9ga724e4h", // StakingCabin
		"sunrise1vwqmjft3svvlqcpudm6wzr3mlvclx5h4wettvt", // Restake
		"sunrise1wnm8t2v4alc067g0erg4uxq3vu7nvpwzfqy7ca", // Stake&Relax
		"sunrise1c5rsy9jqjz6yudm3zhdvpfhlrwly6nh9p4s7yh", // InfraSingularity
		"sunrise18m0fmwesn6za4jefzsfmmv7gx94f38ascc2040", // bountyblok
	}
	// same amount as older validator's one
	coin := sdk.NewInt64Coin("uvrise", 9000000000000)

	fromAddr, err := sdk.AccAddressFromBech32(fromAddress)
	if err != nil {
		panic(err)
	}

	for index, toAddress := range toAddresses {
		toAddr, err := sdk.AccAddressFromBech32(toAddress)
		if err != nil {
			panic(err)
		}
		// if the account is not existent, this method creates account internally
		if err := bankkeeper.SendCoins(ctx, fromAddr, toAddr, sdk.NewCoins(coin)); err != nil {
			panic(err)
		}
		ctx.Logger().Info(fmt.Sprintf("send coin [%s] : target [%s]", strconv.Itoa(index), toAddress))

	}
	return nil
}
