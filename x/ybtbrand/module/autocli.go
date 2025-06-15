package ybtbrand

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
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
					RpcMethod:      "Create",
					Use:            "create [admin]",
					Short:          "Send a create tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "admin"}},
				},
				{
					RpcMethod:      "Mint",
					Use:            "mint [token-creator] [amount]",
					Short:          "Send a mint tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "token_creator"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "Burn",
					Use:            "burn [token-creator] [amount]",
					Short:          "Send a burn tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "token_creator"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "AddYields",
					Use:            "add-yields [token-creator] [amount]",
					Short:          "Send a add-yields tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "token_creator"}, {ProtoField: "amount", Varargs: true}},
				},
				{
					RpcMethod:      "ClaimYields",
					Use:            "claim-yields [token-creator]",
					Short:          "Send a claim-yields tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "token_creator"}},
				},
				{
					RpcMethod:      "UpdateAdmin",
					Use:            "update-admin [new-admin]",
					Short:          "Send a update-admin tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "new_admin"}},
				},
				{
					RpcMethod:      "ClaimCollateralYield",
					Use:            "claim-collateral-yield [token-creator] [base-ybt-creator]",
					Short:          "Send a claim-collateral-yield tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "token_creator"}, {ProtoField: "base_ybt_creator"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
