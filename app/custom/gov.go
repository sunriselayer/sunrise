package custom

import (
	"encoding/json"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

type CustomGovModule struct {
	gov.AppModuleBasic
}

func (cm CustomGovModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := govtypes.DefaultGenesisState()

	day := time.Duration(time.Hour * 24)
	oneWeek := day * 7

	genesis.Params.MinDeposit = sdk.NewCoins(
		sdk.NewCoin("uvrise", math.NewInt(1_000_000_000)),
		sdk.NewCoin("urise", math.NewInt(1_000_000_000*2)),
	)
	genesis.Params.ExpeditedMinDeposit = sdk.NewCoins(
		sdk.NewCoin("uvrise", math.NewInt(1_000_000_000*5)),
		sdk.NewCoin("urise", math.NewInt(1_000_000_000*5*2)),
	)
	genesis.Params.MaxDepositPeriod = &oneWeek
	genesis.Params.VotingPeriod = &oneWeek

	return cdc.MustMarshalJSON(genesis)
}
