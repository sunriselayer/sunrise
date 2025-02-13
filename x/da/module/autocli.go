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
				{
					RpcMethod:      "ProofDeputy",
					Use:            "proof-deputy <validator_address>",
					Short:          "Shows proof deputy of the validator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "validator_address"}},
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
					RpcMethod:      "SubmitInvalidity",
					Use:            "submit-invalidity <metadata_uri> <index>,<index>...",
					Short:          "Submit invalidity to the data",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "metadata_uri"}, {ProtoField: "indices"}},
				},
				{
					RpcMethod: "SubmitValidityProof",
					Skip:      true,
				},
				{
					RpcMethod:      "RegisterProofDeputy",
					Use:            "register-proof-deputy <deputy_address>",
					Short:          "Register proof deputy",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "deputy_address"}},
				},
				{
					RpcMethod:      "UnregisterProofDeputy",
					Use:            "unregister-proof-deputy",
					Short:          "Unregister proof deputy",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
