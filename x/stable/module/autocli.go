package stable

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/sunriselayer/sunrise/x/stable/types"
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
					RpcMethod: "Mint",
					Use:       "mint [amount]",
					Short:     "Send a mint tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "Burn",
					Use:       "burn [amount]",
					Short:     "Send a burn tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount", Varargs: true},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
