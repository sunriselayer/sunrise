# Liquidity Incentive Module

The `x/liquidityincentive` module distributes inflation-based rewards to liquidity providers through a gauge voting system. Users vote on which liquidity pools receive incentives, and can also be influenced by bribes.

## Core Concepts

### Epochs

The module operates in fixed-time intervals called **epochs**. At the start of each epoch, rewards are distributed based on the voting results from the previous epoch.

### Gauges and Voting

Each liquidity pool is associated with a **gauge**. Users vote on these gauges to direct the flow of inflation rewards. A user's voting power is proportional to their staked assets.

### Bribes

To influence voting, users can attach **bribes** (in the form of tokens) to a gauge for a specific epoch. These bribes are then distributed to the voters of that gauge, proportional to their voting power.

### Vote Tallying

At the end of each epoch, the votes for each gauge are tallied. The total inflation rewards for the epoch are then distributed to the liquidity pools according to the proportion of votes each gauge received.

## Messages

### MsgVoteGauge

Casts a vote for one or more liquidity pool gauges.

```protobuf
message MsgVoteGauge {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  repeated PoolWeight pool_weights = 2;
}
```

### MsgRegisterBribe

Registers a bribe for a specific gauge in a future epoch.

```protobuf
message MsgRegisterBribe {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  uint64 epoch_id = 2;
  uint64 pool_id = 3;
  repeated cosmos.base.v1beta1.Coin amount = 4;
}
```

### MsgClaimBribes

Claims earned bribes.

```protobuf
message MsgClaimBribes {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  repeated uint64 bribe_ids = 2;
}
```

### MsgStartNewEpoch

Manually triggers the start of a new epoch. This is a privileged message.

```protobuf
message MsgStartNewEpoch {
  option (cosmos.msg.v1.signer) = "authority";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
}
```

### MsgUpdateParams

Updates the module parameters. This is a governance-gated action.

```protobuf
message MsgUpdateParams {
  option (cosmos.msg.v1.signer) = "authority";
  string authority = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  Params params = 2 [(gogoproto.nullable) = false];
}
```

## Parameters

The `x/liquidityincentive` module has the following parameters:

- `epoch_blocks`: The duration of each epoch in blocks.
- `staking_reward_ratio`: The proportion of inflation rewards allocated to stakers versus liquidity providers.
- `bribe_claim_epochs`: The number of epochs after which unclaimed bribes expire.
