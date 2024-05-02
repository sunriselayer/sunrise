package liquiditypool

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/sunriselayer/sunrise-app/api/sunrise/liquiditypool/v1"
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
					RpcMethod: "PairAll",
					Use:       "list-pair",
					Short:     "List all pair",
				},
				{
					RpcMethod:      "Pair",
					Use:            "show-pair [base_denom] [quote_denom]",
					Short:          "Shows a pair",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "base_denom"}, {ProtoField: "quote_denom"}},
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
					RpcMethod: "TwapAll",
					Use:       "list-twap",
					Short:     "List all twap",
				},
				{
					RpcMethod:      "Twap",
					Use:            "show-twap [base_denom] [quote_denom]",
					Short:          "Shows a twap",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "base_denom"}, {ProtoField: "quote_denom"}},
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
					Use:            "create-pool [base_denom] [quote_denom]",
					Short:          "Create pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "base_denom"}, {ProtoField: "quote_denom"}},
				},
				{
					RpcMethod:      "UpdatePool",
					Use:            "update-pool [id] [base_denom] [quote_denom]",
					Short:          "Update pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "base_denom"}, {ProtoField: "quote_denom"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
