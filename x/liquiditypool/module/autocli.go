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
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "lowerTick"}, {ProtoField: "upperTick"}},
				},
				{
					RpcMethod:      "UpdatePool",
					Use:            "update-pool [id] [lowerTick] [upperTick]",
					Short:          "Update pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "lowerTick"}, {ProtoField: "upperTick"}},
				},
				{
					RpcMethod:      "DeletePool",
					Use:            "delete-pool [id]",
					Short:          "Delete pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
