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
				{
					RpcMethod:      "CreatePosition",
					Use:            "create-position ",
					Short:          "Create position",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod:      "UpdatePosition",
					Use:            "update-position [id] ",
					Short:          "Update position",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "DeletePosition",
					Use:            "delete-position [id]",
					Short:          "Delete position",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "CollectFees",
					Use:            "collect-fees",
					Short:          "Send a collect-fees tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod:      "CollectIncentives",
					Use:            "collect-incentives",
					Short:          "Send a collect-incentives tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
