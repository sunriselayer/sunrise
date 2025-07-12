# Swap Module

The `x/swap` module provides functionality for swapping tokens on the Sunrise blockchain. It supports complex swap routes and can be used as a middleware for IBC transfers.

## Core Concepts

### Routes

The swap module uses a flexible routing system that allows for complex swaps involving multiple pools. A route can be a single pool, a series of pools, or a parallel set of pools.

```protobuf
message Route {
  string denom_in = 1;
  string denom_out = 2;
  oneof strategy {
    RoutePool pool = 3;
    RouteSeries series = 4;
    RouteParallel parallel = 5;
  }
}

message RoutePool {
  uint64 pool_id = 1;
}

message RouteSeries {
  repeated Route routes = 1 [(gogoproto.nullable) = false];
}

message RouteParallel {
  repeated Route routes = 1 [(gogoproto.nullable) = false];
  repeated string weights = 2 [(cosmos_proto.scalar) = "cosmos.Dec"];
}
```

### Interface Fee

The swap module supports an interface fee, which can be charged by front-end applications that facilitate swaps. The fee is a percentage of the swapped amount and is paid to the interface provider.

## Messages

### MsgSwapExactAmountIn

Swaps a specific amount of input tokens for a minimum amount of output tokens.

```protobuf
message MsgSwapExactAmountIn {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string interface_provider = 2 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  Route route = 3;
  cosmossdk.io.math.Int amount_in = 4;
  cosmossdk.io.math.Int min_amount_out = 5;
}
```

### MsgSwapExactAmountOut

Swaps a maximum amount of input tokens for a specific amount of output tokens.

```protobuf
message MsgSwapExactAmountOut {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string interface_provider = 2 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  Route route = 3;
  cosmossdk.io.math.Int max_amount_in = 4;
  cosmossdk.io.math.Int amount_out = 5;
}
```

## IBC Middleware

The swap module can be used as a middleware for IBC transfers, allowing users to swap tokens as part of a cross-chain transfer. This is done by including a `swap` field in the `memo` of the IBC transfer packet.

### Metadata

The `memo` field should contain a JSON object with a `swap` field, which has the following structure:

```typescript
type SwapMetadata = {
  interface_provider: string;
  route: Route;

  forward?: ForwardMetadata;
} & (
  | {
      exact_amount_in: {
        min_amount_out: string;
      };
    }
  | {
      exact_amount_out: {
        amount_out: string;
        change?: ForwardMetadata;
      };
    }
);
```

## Parameters

The `x/swap` module has the following parameters:

- `interface_fee_rate`: The fee rate charged by interface providers.