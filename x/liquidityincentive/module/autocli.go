package liquidityincentive

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/sunriselayer/sunrise/api/sunrise/liquidityincentive"
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
					RpcMethod: "Epochs",
					Use:       "list-epoch",
					Short:     "List all epoch",
				},
				{
					RpcMethod:      "Epoch",
					Use:            "show-epoch [id]",
					Short:          "Shows a epoch by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "Gauges",
					Use:            "list-gauge [previous_epoch_id]",
					Short:          "List all gauge",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "previous_epoch_id"}},
				},
				{
					RpcMethod:      "Gauge",
					Use:            "show-gauge [previous_epoch_id] [pool_id]",
					Short:          "Shows a gauge",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "previous_epoch_id"}, {ProtoField: "pool_id"}},
				},
				{
					RpcMethod:      "PositionIncentives",
					Use:            "incentives-by-position-ids [ids]",
					Short:          "Show incentives by position ids",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ids"}},
				},
				{
					RpcMethod:      "PositionIncentives",
					Use:            "incentives-by-position-id [id]",
					Short:          "Show incentives by position id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "AddressIncentives",
					Use:            "incentives-by-address [address]",
					Short:          "Show incentives by position id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
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
					RpcMethod:      "VoteGauge",
					Use:            "vote-gauge",
					Short:          "Send a vote-gauge tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod:      "CollectIncentiveRewards",
					Use:            "collect-incentive-rewards",
					Short:          "Send a collect-incentive-rewards tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod:      "CollectVoteRewards",
					Use:            "collect-vote-rewards",
					Short:          "Send a collect-vote-rewards tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
