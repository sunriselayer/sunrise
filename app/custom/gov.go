package custom

import (
	"encoding/json"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/app/consts"

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
		sdk.NewCoin(consts.BondDenom, math.NewInt(100_000_000)),
	)
	genesis.Params.ExpeditedMinDeposit = sdk.NewCoins(
		sdk.NewCoin(consts.BondDenom, math.NewInt(20_000_000)),
	)
	genesis.Params.MaxDepositPeriod = &oneWeek
	genesis.Params.VotingPeriod = &oneWeek

	// 20.0%
	// Because we disallow validators to vote with non voting delegators' power,
	// we need to allow lower quorum.
	genesis.Params.Quorum = math.LegacyNewDecWithPrec(200, 3).String()

	return cdc.MustMarshalJSON(genesis)
}
