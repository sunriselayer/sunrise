package lending

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/sunriselayer/sunrise/x/lending/types"
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
					RpcMethod:      "Supply",
					Use:            "supply [amount]",
					Short:          "Send a supply tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}},
				},
				{
					RpcMethod:      "Borrow",
					Use:            "borrow [borrow-denom] [collateral-pool-id] [collateral-position-id]",
					Short:          "Send a borrow tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "borrow_denom"}, {ProtoField: "collateral_pool_id"}, {ProtoField: "collateral_position_id"}},
				},
				{
					RpcMethod:      "Repay",
					Use:            "repay [borrow-id] [amount]",
					Short:          "Send a repay tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "borrow_id"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "Liquidate",
					Use:            "liquidate [borrow-id] [amount]",
					Short:          "Send a liquidate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "borrow_id"}, {ProtoField: "amount"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
