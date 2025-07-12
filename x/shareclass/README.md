# Share Class Module

The `x/shareclass` module provides a non-voting delegation mechanism for the Sunrise blockchain. It allows users to delegate tokens to validators without receiving voting power, while still earning rewards through a share token system.

## Core Concepts

### Non-Voting Delegation

Users can delegate their tokens to validators without receiving voting power. This is useful for users who want to earn staking rewards without participating in governance.

### Share Tokens

When a user delegates tokens, they receive non-transferable share tokens in return. These tokens represent the user's share of the validator's total delegation. The number of share tokens received is proportional to the amount of tokens delegated.

When a user undelegates, their share tokens are burned, and they receive their original tokens back, plus any accumulated rewards.

### Rewards Distribution

Rewards are distributed to delegators based on their share of the validator's total delegation. The rewards are calculated and distributed periodically, as defined by the `reward_period` parameter.

## Messages

### MsgCreateValidator

Creates a new validator.

```protobuf
message MsgCreateValidator {
  option (cosmos.msg.v1.signer) = "sender";
  string validator_address = 1 [(cosmos_proto.scalar) = "cosmos.ValidatorAddressString"];
  cosmos.staking.v1beta1.Description description = 2;
  cosmos.staking.v1beta1.CommissionRates commission = 3;
  cosmossdk.io.math.Int min_self_delegation = 4;
  google.protobuf.Any pubkey = 5;
  cosmos.base.v1beta1.Coin amount = 6;
}
```

### MsgNonVotingDelegate

Delegates tokens to a validator without receiving voting power.

```protobuf
message MsgNonVotingDelegate {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string validator_address = 2 [(cosmos_proto.scalar) = "cosmos.ValidatorAddressString"];
  cosmos.base.v1beta1.Coin amount = 3;
}
```

### MsgNonVotingUndelegate

Undelegates tokens from a validator.

```protobuf
message MsgNonVotingUndelegate {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string validator_address = 2 [(cosmos_proto.scalar) = "cosmos.ValidatorAddressString"];
  cosmos.base.v1beta1.Coin amount = 3;
  string recipient = 4 [(cosmos_proto.scalar) = "cosmos.AddressString"];
}
```

### MsgClaimRewards

Claims the rewards earned from a validator.

```protobuf
message MsgClaimRewards {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string validator_address = 2 [(cosmos_proto.scalar) = "cosmos.ValidatorAddressString"];
}
```

## Parameters

The `x/shareclass` module has the following parameters:

- `reward_period`: The period of time after which rewards are distributed.
- `create_validator_gas`: The amount of gas consumed when creating a new validator.
