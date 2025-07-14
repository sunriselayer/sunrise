# Stable Module

The `x/stable` module is responsible for minting and burning stablecoins.

## Core Concepts

The stable module provides a simple mechanism for managing the supply of a stablecoin. The module has two primary functions:

- **Minting**: The module can mint new stablecoins. This is a privileged operation that can only be performed by an authority address.
- **Burning**: The module can burn existing stablecoins. This is also a privileged operation that can only be performed by an authority address.

## Messages

### MsgMint

Mints new stablecoins.

```protobuf
message MsgMint {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  cosmossdk.io.math.Int amount = 2;
}
```

### MsgBurn

Burns existing stablecoins.

```protobuf
message MsgBurn {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  repeated cosmos.base.v1beta1.Coin amount = 2;
}
```

## Parameters

The `x/stable` module has the following parameters:

- `stable_denom`: The denomination of the stablecoin.
- `authority_addresses`: A list of addresses that are authorized to mint and burn stablecoins.
