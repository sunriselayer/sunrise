package liquiditypool

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "Pools",
					Use:       "list-pools",
					Short:     "List all pools",
				},
				{
					RpcMethod:      "Pool",
					Use:            "show-pool <id>",
					Short:          "Shows a pool by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "Positions",
					Use:       "list-positions",
					Short:     "List all positions",
				},
				{
					RpcMethod:      "Position",
					Use:            "show-position <id>",
					Short:          "Shows a position by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "PoolPositions",
					Use:            "show-pool-positions <pool_id>",
					Short:          "List positions by pool id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}},
				},
				{
					RpcMethod:      "AddressPositions",
					Use:            "address-positions <address>",
					Short:          "List positions by address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod:      "PositionFees",
					Use:            "show-position-fees <id>",
					Short:          "Shows fees in a position by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "CreatePool",
					Use:       "create-pool <denom_base> <denom_quote> <fee_rate> <price_ratio> <base_offset>",
					Short:     "Create pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "denom_base"},
						{ProtoField: "denom_quote"},
						{ProtoField: "fee_rate"},
						{ProtoField: "price_ratio"},
						{ProtoField: "base_offset"},
					},
				},
				{
					RpcMethod: "CreatePosition",
					Use:       "create-position <pool_id> <lower_tick> <upper_tick> <token_base.amount> <token_base.denom> <token_quote.amount> <token_quote.denom> <min_amount_base> <min_amount_quote>",
					Short:     "Create position",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "pool_id"},
						{ProtoField: "lower_tick"},
						{ProtoField: "upper_tick"},
						{ProtoField: "token_base.amount"},
						{ProtoField: "token_base.denom"},
						{ProtoField: "token_quote.amount"},
						{ProtoField: "token_quote.denom"},
						{ProtoField: "min_amount_base"},
						{ProtoField: "min_amount_quote"},
					},
				},
				{
					RpcMethod: "IncreaseLiquidity",
					Use:       "increase-liquidity <id> <amount_base> <amount_quote> <min_amount_base> <min_amount_quote>",
					Short:     "Increase liquidity",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
						{ProtoField: "amount_base"},
						{ProtoField: "amount_quote"},
						{ProtoField: "min_amount_base"},
						{ProtoField: "min_amount_quote"},
					},
				},
				{
					RpcMethod: "DecreaseLiquidity",
					Use:       "decrease-liquidity <id> <liquidity>",
					Short:     "Decrease liquidity",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
						{ProtoField: "liquidity"},
					},
				},
				{
					RpcMethod: "ClaimRewards",
					Use:       "claim-rewards <position_id>...",
					Short:     "Send a claim-rewards tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "position_ids"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
