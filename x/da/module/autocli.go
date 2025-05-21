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
					Use:            "published-data [metadata_uri]",
					Short:          "Shows published data",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "metadata_uri"}},
				},
				{
					RpcMethod: "AllPublishedData",
					Use:       "all-published-data",
					Short:     "Shows all published data",
				},
				{
					RpcMethod:      "ValidityProof",
					Use:            "validity-proof [metadata_uri] [validator_address]",
					Short:          "Shows the validity proof of the data by validator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "metadata_uri"}, {ProtoField: "validator_address"}},
				},
				{
					RpcMethod:      "AllValidityProofs",
					Use:            "all-validity-proofs [metadata_uri]",
					Short:          "Shows all validity proofs of the data",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "metadata_uri"}},
				},
				{
					RpcMethod:      "Invalidity",
					Use:            "invalidity [metadata_uri] [sender_address]",
					Short:          "Shows invalidity of the data by sender",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "metadata_uri"}, {ProtoField: "sender_address"}},
				},
				{
					RpcMethod:      "AllInvalidity",
					Use:            "all-invalidity [metadata_uri]",
					Short:          "Shows all invalidity of the data",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "metadata_uri"}},
				},
				{
					RpcMethod:      "ZkpProofThreshold",
					Use:            "zkp-proof-threshold [shard_count]",
					Short:          "Shows threshold number of proof",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "shard_count"}},
				},
				{
					RpcMethod:      "ProofDeputy",
					Use:            "proof-deputy [validator_address]",
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
					Use:            "submit-invalidity [metadata_uri] [index],[index],...",
					Short:          "Send a submit-invalidity tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "metadata_uri"}, {ProtoField: "indices"}},
				},
				{
					RpcMethod: "SubmitValidityProof",
					Skip:      true,
				},
				{
					RpcMethod:      "RegisterProofDeputy",
					Use:            "register-proof-deputy [deputy_address]",
					Short:          "Send a register-proof-deputy tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "deputy_address"}},
				},
				{
					RpcMethod:      "UnregisterProofDeputy",
					Use:            "unregister-proof-deputy",
					Short:          "Send a unregister-proof-deputy tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
