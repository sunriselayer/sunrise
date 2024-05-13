# liquiditypool

## Basic architecture

- There are many `Pool`s in the module
  - Each `Pool` has a base denom, quote denom, fee rate and Constant Function Market Maker (CFMM), and a flag whether it is concentrated liquidity pool or not.
  - `$SR` can't be a base denom.
  - The token `$XXX` which already has `$XXX/$SR` pair only can be used for the quote denom.
    - This is for making the price able to convert to the `$SR` based value.
- TWAP are recorded for each pair of base and quote denom after the swapping.
- The LP token denom will be `liquiditypool/{pool_id}` if it is not concentrated liquidity pool.

## Constant Function Market Maker (CFMM)

### Constant Product Market Maker

$$
  x y = k
$$

Concentrated Liquidity is supported.

### Stable Market Maker

$$
  x y (x^2 + y^2) = k
$$

Concentrated Liquidity is not supported.

## Concentrated Liquidity Pool

If the flag whether it is concentrated liquidity pool or not is set to true, the pool is concentrated liquidity pool.

The LP token denom will be `liquiditypool/{pool_id}/range/{range_index}`.

## Params

- `swap_treasury_tax_rate`
  - Some amounts go to the treasury module address `moduleaddress("liquiditypool/{pool_id}/treasury")`
- `swap_interface_fee_rate`
  - The frontend interface provider for each swap tx can receive the fee.
- `pool_exit_fee_rate`
- `twap_window`
- `twap_expiry`

### Validation

- All fee rates must be positive.
- Sum of swap fee rates must be less than 1.
