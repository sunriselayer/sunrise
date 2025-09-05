package types

import (
	"encoding/json"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestPrintAminoMsgSwapExactAmountIn(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	RegisterLegacyAminoCodec(cdc)

	msg := &MsgSwapExactAmountIn{
		Sender:            "sunrise1xxgjt7yqkmn63m2d0nrf0vt5uuc2hr6l45xaa9",
		InterfaceProvider: "sunrise1xxgjt7yqkmn63m2d0nrf0vt5uuc2hr6l45xaa9",
		Route: Route{
			DenomIn:  "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
			DenomOut: "ibc/A7AD825A4B48DDA0138D118655E60100D22A4D690C45B95221520B58C9A64B63",
			Strategy: &Route_Series{
				Series: &RouteSeries{
					Routes: []Route{
						{
							DenomIn:  "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
							DenomOut: "uusdrise",
							Strategy: &Route_Pool{
								Pool: &RoutePool{
									PoolId: 2,
								},
							},
						},
						{
							DenomIn:  "uusdrise",
							DenomOut: "ibc/A7AD825A4B48DDA0138D118655E60100D22A4D690C45B95221520B58C9A64B63",
							Strategy: &Route_Pool{
								Pool: &RoutePool{
									PoolId: 3,
								},
							},
						},
					},
				},
			},
		},
		AmountIn:     math.NewInt(1000000),
		MinAmountOut: math.NewInt(982183),
	}

	aminoJSON, err := cdc.MarshalJSON(msg)
	require.NoError(t, err)

	var prettyJSON map[string]interface{}
	err = json.Unmarshal(aminoJSON, &prettyJSON)
	require.NoError(t, err)

	prettyPrinted, err := json.MarshalIndent(prettyJSON, "", "  ")
	require.NoError(t, err)

	fmt.Println(string(prettyPrinted))
}
