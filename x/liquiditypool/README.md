# liquiditypool

The module `x/liquiditypool` is for liquidity poll with concentrated liquidity AMM mechanism.

## Spec

### Pool

Each pool has these parameters

- `id`
- `denom_base`
- `denom_quote`
- `fee_rate`
- `tick_params`
- `current_tick`
- `current_tick_liquidity`
- `current_sqrt_price`

### Tick

There are two parameters in each pool

- `price_ratio`: typically `1.0001`
- `base_offset`: typically `0` otherwise in `[0, 1)`

The tick-price conversion formula is this:

$$
\text{price}(\text{tick}) = \text{price\_ratio} ^ {\text{tick} - \text{base\_offset}}
$$

In the typiical case,

$$
\text{price}(\text{tick}) = 1.0001 ^ {\text{tick}}
$$

and it is same to Uniswap V3.

### Position

Each position has these info

- `id`: Unique ID
- `address`: Sender's address
- `pool_id`: Pool's ID
- `lower_tick`: Uniquely determine the min price of the range of the position
- `upper_tick`: Uniquely determine the max price of the range of the position
- `liquidity`: The amount of liquidity determined based on the volume of deposits

## Swap

For sending the tx for swapping tokens, use msgs in `x/swap` module.

## Messages

### MsgCreatePool

Users can create a Liquidity Pool.

- `authority`: address with pool authority
- `denom_base`: base denom
- `denom_quote`: quote denom
- `fee_rate`: fee rate charged on swaps
- `price_ratio`: price-ratio parameter (default 1.0001)
- `base_offset`: base-offset parameter (default 0.5)

### MsgCreatePosition

Users can create a position for the pool.

- `sender`: sender address
- `pool_id`: ID of the liquidity pool
- `lower_tick`: tick of lower price
- `upper_tick`: tick of upper price
- `token_base`: base denom & amount (cosmos/base/v1bata1/coin type)
- `token_quote`: quote denom & amount (cosmos/base/v1bata1/coin type)
- `min_amount_base`: Minimum base amount
- `min_amount_quote`: Minimum quote amount

For example, if `lower_tick` is -4155 & `upper_tick` is 4054, the price range is -0.66 ~ 1.5.

if the actual value falls below either `min_amount_base` or `min_amount_quote`, the position creation will be canceled. To avoid this check, set them to 0.

### MsgIncreaseLiquidity

Users can increase liquidity of an existing position.

- `sender`: sender address
- `pool_id`: ID of the liquidity pool
- `amount_base`: base amount to add
- `amount_quote`: quote amount to add
- `min_amount_base`: Minimum base amount
- `min_amount_quote`: Minimum quote amount

`lower_tick` and `upper_tick` are not changed from the existing position.

### MsgDecreaseLiquidity

Users can decrease liquidity of an existing position.

- `sender`: sender address
- `id`: ID of the liquidity pool
- `liquidity`: Amount of liquidity to be decreased

The `liquidity` value of the current position can be gotten by a query.
If the value equal to or greater than the existing liquidity is set, the position will be deleted.

### MsgClaimRewards

Users can claim fees & incentives for providing liquidity.
Fees are `base_denom` or `quote_denom`, incentives are provided by vRISE.

- `position_ids`: The list of position ids to claim rewards

## Query

See [openapi.yml](../../docs/static/openapi.yml) for details

- Params
- Pools
- Pool
- Positions
- Position
- PoolPositions
- AddressPositions
- PositionFees
- CalculationCreatePosition
- CalculationIncreaseLiquidity
