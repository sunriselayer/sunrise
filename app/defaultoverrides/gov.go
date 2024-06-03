package defaultoverrides

import (
	"encoding/json"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	gov "github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// GovModuleBasic is a custom wrapper around the x/gov module's AppModuleBasic
// implementation to provide custom default genesis state.
type GovModuleBasic struct {
	gov.AppModuleBasic
}

// DefaultGenesis returns custom x/gov module genesis state.
func (GovModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genState := govtypes.DefaultGenesisState()
	day := time.Duration(time.Hour * 24)
	oneWeek := day * 7

	genState.Params.MinDeposit = sdk.NewCoins(
		sdk.NewCoin("uvrise", sdkmath.NewInt(1_000_000_000)),
		sdk.NewCoin("urise", sdkmath.NewInt(1_000_000_000*2)),
	)
	genState.Params.ExpeditedMinDeposit = sdk.NewCoins(
		sdk.NewCoin("uvrise", sdkmath.NewInt(1_000_000_000*5)),
		sdk.NewCoin("urise", sdkmath.NewInt(1_000_000_000*5*2)),
	)
	genState.Params.MaxDepositPeriod = &oneWeek
	genState.Params.VotingPeriod = &oneWeek

	return cdc.MustMarshalJSON(genState)
}
