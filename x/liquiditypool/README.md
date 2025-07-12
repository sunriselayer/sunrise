# Liquidity Pool Module

The `x/liquiditypool` module implements a concentrated liquidity automated market maker (AMM), based on the concepts introduced by Uniswap V3.

## Core Concepts

### Concentrated Liquidity

Unlike traditional AMMs where liquidity is distributed uniformly along the entire price curve, concentrated liquidity allows liquidity providers (LPs) to allocate their capital to specific price ranges. This results in higher capital efficiency, as LPs can earn more fees on the same amount of capital.

### Ticks

Ticks are discrete price points that define the boundaries of a price range. The price of an asset is determined by the current tick of the pool. The relationship between ticks and prices is defined by the following formula:

```
price(tick) = price_ratio ^ (tick - base_offset)
```

- `price_ratio`: The price ratio between two consecutive ticks. Typically set to `1.0001`.
- `base_offset`: An offset to the tick, allowing for more granular control over the price range. Typically set to `0`.

### Positions

LPs provide liquidity by creating positions, which are defined by a lower and upper tick. The liquidity provided by a position is only active when the current price of the pool is within the range defined by the position's ticks.

## Messages

### MsgCreatePool

Creates a new liquidity pool.

```protobuf
message MsgCreatePool {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  string denom_base = 2;
  string denom_quote = 3;
  string fee_rate = 4;
  string price_ratio = 5;
  string base_offset = 6;
}
```

### MsgCreatePosition

Creates a new position in a liquidity pool.

```protobuf
message MsgCreatePosition {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  uint64 pool_id = 2;
  int64 lower_tick = 3;
  int64 upper_tick = 4;
  cosmos.base.v1beta1.Coin token_base = 5;
  cosmos.base.v1beta1.Coin token_quote = 6;
  cosmossdk.io.math.Int min_amount_base = 7;
  cosmossdk.io.math.Int min_amount_quote = 8;
}
```

### MsgIncreaseLiquidity

Increases the liquidity of an existing position.

```protobuf
message MsgIncreaseLiquidity {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  uint64 id = 2;
  cosmossdk.io.math.Int amount_base = 3;
  cosmossdk.io.math.Int amount_quote = 4;
  cosmossdk.io.math.Int min_amount_base = 5;
  cosmossdk.io.math.Int min_amount_quote = 6;
}
```

### MsgDecreaseLiquidity

Decreases the liquidity of an existing position.

```protobuf
message MsgDecreaseLiquidity {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  uint64 id = 2;
  string liquidity = 3;
}
```

### MsgClaimRewards

Claims the fees and incentives earned by a position.

```protobuf
message MsgClaimRewards {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  repeated uint64 position_ids = 2;
}
```

## Parameters

The `x/liquiditypool` module has the following parameters:

- `create_pool_gas`: The amount of gas consumed when creating a new pool.
- `withdraw_fee_rate`: The fee rate charged when withdrawing liquidity. (Currently not used)
- `allowed_quote_denoms`: A list of denoms that are allowed to be used as the quote asset in a pool.