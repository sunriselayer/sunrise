# Data Availability (DA) Module

The `x/da` module is responsible for ensuring data availability on the Sunrise blockchain. It uses a combination of erasure coding, zero-knowledge proofs (ZKPs), and a challenge mechanism to verify that data is being stored correctly by validators.

## Core Concepts

### Data Publication

L2 chains can publish data to the DA layer by sending a `MsgPublishData` transaction. This data is split into shards, which are then erasure coded to provide redundancy. The double hashes of the shards are stored on-chain.

### Zero-Knowledge Proofs

ZKPs are used to verify that a validator is in possession of the data shards they are responsible for, without revealing the data itself. This is done by verifying a proof against the double hashes of the shards.

### Challenge Mechanism

After data is published, there is a challenge period during which anyone can challenge the validity of the data by submitting a `MsgSubmitInvalidity` transaction. If a challenge is successful, the data is marked as invalid. If the challenge period expires without any successful challenges, the data is marked as verified.

### Slashing

Validators are responsible for voting on the validity of published data. Validators who fail to vote or who vote fraudulently are slashed.

## Messages

### MsgPublishData

Publishes data to the DA layer.

```protobuf
message MsgPublishData {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string metadata_uri = 2;
  uint64 parity_shard_count = 3;
  repeated bytes shard_double_hashes = 4;
  string data_source_info = 5;
}
```

### MsgSubmitInvalidity

Submits a challenge to the validity of published data.

```protobuf
message MsgSubmitInvalidity {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string metadata_uri = 2;
  repeated int64 indices = 3;
}
```

### MsgSubmitValidityProof

Submits a proof of validity for challenged data.

```protobuf
message MsgSubmitValidityProof {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string validator_address = 2 [(cosmos_proto.scalar) = "cosmos.ValidatorAddressString"];
  string metadata_uri = 3;
  repeated int64 indices = 4;
  repeated bytes proofs = 5;
}
```

### MsgRegisterProofDeputy

Registers a deputy account to submit validity proofs on behalf of a validator.

```protobuf
message MsgRegisterProofDeputy {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string deputy_address = 2 [(cosmos_proto.scalar) = "cosmos.AddressString"];
}
```

### MsgUnregisterProofDeputy

Unregisters a deputy account.

```protobuf
message MsgUnregisterProofDeputy {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
}
```

### MsgVerifyData

Triggers the verification of data.

```protobuf
message MsgVerifyData {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
}
```

## Parameters

The `x/da` module has the following parameters:

- `publish_data_gas`: The amount of gas consumed when publishing data.
- `challenge_threshold`: The threshold of invalid shards required to challenge data.
- `replication_factor`: The number of validators that should store a copy of the data.
- `slash_epoch`: The number of blocks in a slash epoch.
- `slash_fault_threshold`: The threshold of faults a validator can have in an epoch before being slashed.
- `slash_fraction`: The fraction of a validator's stake that is slashed.
- `challenge_period`: The duration of the challenge period.
- `proof_period`: The duration of the proof period.
- `rejected_removal_period`: The duration after which rejected data is removed.
- `verified_removal_period`: The duration after which verified data is removed.
- `publish_data_collateral`: The collateral required to publish data.
- `submit_invalidity_collateral`: The collateral required to submit an invalidity challenge.
- `zkp_verifying_key`: The ZKP verifying key.
- `zkp_proving_key`: The ZKP proving key.
- `min_shard_count`: The minimum number of shards.
- `max_shard_count`: The maximum number of shards.
- `max_shard_size`: The maximum size of a shard.
