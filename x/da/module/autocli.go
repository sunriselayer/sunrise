package da

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/sunriselayer/sunrise/x/da/types"
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
					RpcMethod:      "PublishedData",
					Use:            "published-data <metadata_uri>",
					Short:          "Shows published data",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "metadata_uri"}},
				},
				{
					RpcMethod: "AllPublishedData",
					Use:       "all-published-data",
					Short:     "Shows all published data",
				},
				{
					RpcMethod:      "ZkpProofThreshold",
					Use:            "zkp-proof-threshold <shard_count>",
					Short:          "Shows threshold number of proof",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "shard_count"}},
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
					RpcMethod: "PublishData",
					Skip:      true,
				},
				{
					RpcMethod: "SubmitInvalidity",
					Skip:      true,
				},
				{
					RpcMethod: "SubmitValidityProof",
					Skip:      true,
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
