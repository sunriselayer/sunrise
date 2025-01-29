# DA

## Specification for Zero-Knowledge Proof

### Terms and Notation

- $n$: Total number of erasure coded shards
- $t$: Threshold (minimum number of shards required to prove the possession)
- $s_i$: The $i$-th erasure coded data shard
- $H$: The hash function

### Overview

This system verifies the possession of data shard hash `H(s_i)` without exposing `H(s_i)`.

### Zero-Knowledge Proof System

#### Public Inputs

- $\{H^2(s_i)\}_{i=1}^n$
- $t$: Threshold

#### Private Inputs

- $I$: Index set of shards
- $\{H(s_i)\}_{i \in I}$

#### ZKP Circuit Constraints

- For each $i$, verify $H^2(s_i) = H^2(s_i)_{public}$
- $t \le |I|$

## Msgs

### MsgPublishData

```protobuf
message MsgPublishData {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1;
  string metadata_uri = 2;
  repeated bytes shard_double_hashes = 3;
}
```

## Metadata

The metadata that is preserved as a JSON in the `metadata_uri` has a schema below

```protobuf
message Metadata {
  uint64 shard_size = 1;
  uint64 shard_count = 2;
  repeated string shard_uris = 3;
}
```

## URI

Supported URIs are below

- `ipfs://`: IPFS
- `ar://`: Arweave

## Status

PublishedData sent from the L2 chain to the DA layer changes status depending on the transmission of Tx and the passage of time.

```protobuf
enum Status {
  // Default value
  STATUS_UNSPECIFIED = 0;
  // Verified
  STATUS_VERIFIED = 1;
  // Rejected
  STATUS_REJECTED = 2;
  // After processing in the msg_server
  STATUS_VOTING = 3;
  // Verified the votes from the validators. Challenge can be received (after preBlocker)
  STATUS_CHALLENGE_PERIOD = 4;
  // reported as fraud. SubmitProof tx can be received (after received ChallengeForFraud tx)
  STATUS_CHALLENGING = 5;
}
```

1. A L2 chain sends the MsgPublishData transaction via sunrise-data, etc. If Tx is successful, it is registered with `VOTING` status.
1. Registered PublishedData will be changed to `CHALLENGE_PERIOD` status in PreBlocker.
1. During `CHALLENGE_PERIOD`, the status can be changed to `CHALLENGING` status through MsgChallengeForFraud Tx by anyone.
1. In EndBlocker, `CHALLENGE_PERIOD` that has passed ChallengePeriod become `VERIFIED`. `CHALLENGING` will become `REJECTED` if `valid_shards < data_shard_count`. Otherwise, it will become `VERIFIED`.
