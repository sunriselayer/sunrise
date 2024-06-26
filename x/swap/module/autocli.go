package swap

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/sunriselayer/sunrise/api/sunrise/swap"
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
					RpcMethod: "OutgoingInFlightPackets",
					Use:       "list-in-flight-packet",
					Short:     "List all out-going-in-flight-packet",
				},
				{
					RpcMethod:      "OutgoingInFlightPacket",
					Use:            "show-outgoing-in-flight-packet [src-port] [src-channel] [sequence]",
					Short:          "Shows an outgoing-in-flight-packet",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "src_port_id"}, {ProtoField: "src_channel_id"}, {ProtoField: "sequence"}},
				},
				{
					RpcMethod: "IncomingInFlightPackets",
					Use:       "list-incoming-in-flight-packet-packet",
					Short:     "List all incoming-in-flight-packet",
				},
				{
					RpcMethod:      "IncomingInFlightPacket",
					Use:            "incoming-in-flight-packet [src-port] [src-channel] [sequence]",
					Short:          "Shows an incoming-in-flight-packet",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "src_port_id"}, {ProtoField: "src_channel_id"}, {ProtoField: "sequence"}},
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
					RpcMethod:      "SwapExactAmountIn",
					Use:            "swap-exact-amount-in",
					Short:          "Send a swap-exact-amount-in tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod:      "SwapExactAmountOut",
					Use:            "swap-exact-amount-out",
					Short:          "Send a swap-exact-amount-out tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						// {ProtoField: "interface_provider"},
						// {ProtoField: "route"},
						// {ProtoField: "max_amount_in"},
						// {ProtoField: "amount_out"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
