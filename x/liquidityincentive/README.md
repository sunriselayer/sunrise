# liquidityincentive

## Spec

### BeginBlocker

- Create a new `Epoch` if the last `Epoch` has ended or the first `Epoch` has not been created.

### EndBlocker

- Transfer a portion of inflation rewards from `x/distribution` pool to `x/liquidityincentive` pool.

### Lazy accounting

This module utilize Fee Accumulator in `x/liquiditypool` for each position in the liquidity pool to keep track of the rewards accumulated for each position.
To minimize the load of reward distribution, the rewards are not distributed to the positions immediately. Instead, the rewards amount is calculated when users claim the rewards.

For pool $i$ position $j$,

- $\text{PositionUnclaimedAccumulation}_{ij}$ is the accumulated quote-valued fee rewards for the position.
- $\text{Incentive} \leftarrow \Delta t \times \text{InflationRate} \times (1 - \text{ValidatorRewardRatio}) \times \text{IncentiveTotalSupply}$
- $\text{PoolInflated}_i = \text{GaugeWeight}_i \times \text{Incentive}$
- $\text{PoolUnclaimed}_i = \text{PoolInflated}_i - \text{PoolTotalClaimed}_i$
- $\text{PoolUnclaimedAccumulation}_i = \sum_j \text{PositionUnclaimedAccumulation}_{ij}$
- $\text{ClaimAmount}_{ij} = \frac{\text{PositionUnclaimedAccumulation}_{ij}}{\text{PoolUnclaimedAccumulation}_i} \times \text{PoolUnclaimed}_i$
- $\text{PositionUnclaimedAccumulation}_{ij} \leftarrow 0$
- $\text{PoolTotalClaimed}_i \leftarrow \text{PoolTotalClaimed}_i - \text{ClaimAmount}_{ij}$

### Epoch

Three epochs concurrently exist in the system.

- Past Epoch: The epoch that has ended.
- Current Epoch: The epoch that is ongoing.
- Next Epoch: The epoch that will be started after the current epoch.

### MsgVoteGauge

Users can vote for a gauge for next epoch. `previous_epoch_id` will be the id of current epoch.

- `previous_epoch_id`

### MsgCollectVoteRewards

Users can claim vote rewards for the past epoch.

- `epoch_id`

### MsgCollectIncentiveRewards

- `pool_ids`
