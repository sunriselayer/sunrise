package tokenfactory

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/sunriselayer/sunrise/x/tokenfactory/types"
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
				{
					RpcMethod:      "DenomAuthorityMetadata",
					Use:            "denom-authority-metadata [denom]",
					Short:          "Shows the authority metadata for a specific denom",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod:      "DenomsFromCreator",
					Use:            "denoms-from-creator [address]",
					Short:          "Shows all tokens created by a specific creator address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "creator"}},
				},
				// {
				// 	RpcMethod:      "BeforeSendHookAddress",
				// 	Use:            "before-send-hook-address [denom]",
				// 	Short:          "Shows the address registered for the before send hook",
				// 	PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				// },
				// {
				// 	RpcMethod:      "AllBeforeSendHooksAddresses",
				// 	Use:            "all-before-send-hooks",
				// 	Short:          "Shows all before send hooks registered",
				// 	PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				// },
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
				// this line is used by ignite scaffolding # autocli/tx
				{
					RpcMethod:      "CreateDenom",
					Use:            "create-denom [sub_denom]",
					Short:          "Send a create-denom tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "subdenom"}},
				},
				{
					RpcMethod:      "Mint",
					Use:            "mint [amount] [mint_to_address]",
					Short:          "Send a mint tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "mint_to_address"}},
				},
				{
					RpcMethod:      "Burn",
					Use:            "burn [amount] [burn_from_address]",
					Short:          "Send a burn tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "burn_from_address"}},
				},
				{
					RpcMethod:      "ChangeAdmin",
					Use:            "change-admin [denom] [new_admin]",
					Short:          "Send a change-admin tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}, {ProtoField: "new_admin"}},
				},
				// {
				// 	RpcMethod:      "SetBeforeSendHook",
				// 	Use:            "set-before-send-hook [denom] [cosmwasm_address]",
				// 	Short:          "Send a set-before-send-hook tx",
				// 	PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}, {ProtoField: "cosmwasm_address"}},
				// },
				{
					RpcMethod:      "SetDenomMetadata",
					Use:            "set-denom-metadata [metadata]",
					Short:          "Send a set-denom-metadata tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "metadata"}},
				},
				{
					RpcMethod:      "ForceTransfer",
					Use:            "force-transfer [amount] [transfer_from_address] [transfer_to_address]",
					Short:          "Send a force-transfer tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "transfer_from_address"}, {ProtoField: "transfer_to_address"}},
				},
			},
		},
	}
}
