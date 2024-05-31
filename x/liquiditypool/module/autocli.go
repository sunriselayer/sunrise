package liquiditypool

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/sunriselayer/sunrise/api/sunrise/liquiditypool"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "PoolAll",
					Use:       "list-pool",
					Short:     "List all pool",
				},
				{
					RpcMethod:      "Pool",
					Use:            "show-pool [id]",
					Short:          "Shows a pool by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "PositionAll",
					Use:       "list-position",
					Short:     "List all position",
				},
				{
					RpcMethod:      "Position",
					Use:            "show-position [id]",
					Short:          "Shows a position by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreatePool",
					Use:            "create-pool [lowerTick] [upperTick]",
					Short:          "Create pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						// {ProtoField: "denom_base"},
						// {ProtoField: "denom_quote"},
						// {ProtoField: "fee_rate"},
						// {ProtoField: "tick_params"},
					},
				},
				{
					RpcMethod:      "CreatePosition",
					Use:            "create-position ",
					Short:          "Create position",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						// {ProtoField: "pool_id"},
						// {ProtoField: "lower_tick"},
						// {ProtoField: "upper_tick"},
						// {ProtoField: "token_base"},
						// {ProtoField: "token_quote"},
						// {ProtoField: "min_amount_base"},
						// {ProtoField: "min_amount_quote"},
					},
				},
				{
					RpcMethod:      "IncreaseLiquidity",
					Use:            "increase-liquidity",
					Short:          "Increase liquidity",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						// {ProtoField: "id"},
						// {ProtoField: "amount_base"},
						// {ProtoField: "amount_quote"},
						// {ProtoField: "min_amount_base"},
						// {ProtoField: "min_amount_quote"},
					},
				},
				{
					RpcMethod:      "DecreaseLiquidity",
					Use:            "decrease-liquidity",
					Short:          "Decrease liquidity",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						// {ProtoField: "id"},
						// {ProtoField: "liquidity"},
					},
				},
				{
					RpcMethod:      "CollectFees",
					Use:            "collect-fees",
					Short:          "Send a collect-fees tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						// {ProtoField: "position_ids"},
					},
				},
				{
					RpcMethod:      "CollectIncentives",
					Use:            "collect-incentives",
					Short:          "Send a collect-incentives tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						// {ProtoField: "position_ids"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
