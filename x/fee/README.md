# Fee Module

The `x/fee` module is responsible for burning a portion of transaction fees. This mechanism helps to reduce the total supply of the native token, potentially increasing its value over time.

## Core Concepts

### Fee Burning

The fee burning process is handled by the `Burn` function in the keeper. Here's a step-by-step overview:

1.  **Fee Collection**: Transaction fees are collected in the `fee_denom` specified in the module's parameters.
2.  **Burn Ratio**: A portion of the collected fees, determined by the `burn_ratio`, is designated for burning.
3.  **Swapping (if necessary)**: If the `fee_denom` is different from the `burn_denom`, the collected fees are swapped to the `burn_denom` using the liquidity pool specified by `burn_pool_id`.
4.  **Burning**: The designated amount of `burn_denom` is sent to the module's account and then burned, permanently removing it from circulation.

The entire burn process is atomic. If any step fails (e.g., the swap fails), the entire operation is reverted, ensuring that no funds are lost or stuck.

## Messages

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

The `x/fee` module has the following parameters:

- `fee_denom`: The denomination of the token that is accepted as a transaction fee.
- `burn_denom`: The denomination of the token that is burned.
- `burn_ratio`: The portion of the transaction fee that is burned. This is a decimal value between 0 and 1.
- `burn_pool_id`: The ID of the liquidity pool to use for swapping fees to the `burn_denom`.
- `burn_enabled`: A boolean value that enables or disables the fee burning mechanism.