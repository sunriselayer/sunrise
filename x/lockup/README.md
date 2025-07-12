# Lockup Module

The `x/lockup` module provides functionality for locking up tokens and delegating them to validators without voting power.

## Core Concepts

### Lockup Accounts

A lockup account is a special type of account that holds tokens for a specified period of time. The tokens in a lockup account are subject to a vesting schedule, which determines when they can be withdrawn.

### Non-Voting Delegation

Tokens held in a lockup account can be delegated to validators. This allows users to earn staking rewards on their locked tokens without having to participate in governance.

### Rewards

Rewards earned from non-voting delegation are automatically added to the lockup account.

## Messages

### MsgInitLockupAccount

Initializes a new lockup account.

```protobuf
message MsgInitLockupAccount {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string owner = 2 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  int64 start_time = 3;
  int64 end_time = 4;
  cosmos.base.v1beta1.Coin amount = 5;
}
```

### MsgNonVotingDelegate

Delegates tokens from a lockup account to a validator.

```protobuf
message MsgNonVotingDelegate {
  option (cosmos.msg.v1.signer) = "owner";
  string owner = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  uint64 lockup_account_id = 2;
  string validator_address = 3 [(cosmos_proto.scalar) = "cosmos.ValidatorAddressString"];
  cosmos.base.v1beta1.Coin amount = 4;
}
```

### MsgNonVotingUndelegate

Undelegates tokens from a validator back to a lockup account.

```protobuf
message MsgNonVotingUndelegate {
  option (cosmos.msg.v1.signer) = "owner";
  string owner = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  uint64 lockup_account_id = 2;
  string validator_address = 3 [(cosmos_proto.scalar) = "cosmos.ValidatorAddressString"];
  cosmos.base.v1beta1.Coin amount = 4;
}
```

### MsgClaimRewards

Claims staking rewards and adds them to the lockup account.

```protobuf
message MsgClaimRewards {
  option (cosmos.msg.v1.signer) = "owner";
  string owner = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  uint64 lockup_account_id = 2;
  string validator_address = 3 [(cosmos_proto.scalar) = "cosmos.ValidatorAddressString"];
}
```

### MsgSend

Sends tokens from a lockup account to another account.

```protobuf
message MsgSend {
  option (cosmos.msg.v1.signer) = "owner";
  string owner = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  uint64 lockup_account_id = 2;
  string recipient = 3 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  repeated cosmos.base.v1beta1.Coin amount = 4;
}
```

## Parameters

The `x/lockup` module has the following parameters:

- `min_lockup_duration`: The minimum duration for which tokens can be locked up.
