package selfdelegation

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/sunriselayer/sunrise/x/selfdelegation/types"
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
				{
					RpcMethod:      "SelfDelegationProxyAccountByOwner",
					Use:            "self-delegation-proxy-account-by-owner [owner_address]",
					Short:          "Shows the self-delegation proxy account by owner address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "owner_address"}},
				},
				{
					RpcMethod:      "LockupAccountsByOwner",
					Use:            "lockup-accounts-by-owner [owner_address]",
					Short:          "Shows the lockup accounts by owner address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "owner_address"}},
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
					RpcMethod:      "SelfDelegate",
					Use:            "self-delegate <amount>",
					Short:          "Send a self-delegate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}},
				},
				{
					RpcMethod:      "WithdrawSelfDelegationUnbonded",
					Use:            "withdraw-self-delegation-unbonded <amount>",
					Short:          "Send a withdraw-self-delegation-unbonded tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}},
				},
				{
					RpcMethod: "RegisterLockupAccount",
					Skip:      true,
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
