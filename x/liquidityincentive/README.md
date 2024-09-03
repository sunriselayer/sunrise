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
- $\text{Incentive} \leftarrow \Delta t \times \text{InflationRate} \times (1 - \text{StakingRewardRatio}) \times \text{IncentiveTotalSupply}$
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
<!-- - Next Epoch: The epoch that will be started after the current epoch. -->

Each epoch has these parameters

- `start_block`
- `end_block`
- `gauges`

## Msg

### MsgVoteGauge

Users can vote for the gauge.

- `weights`: The list of what Percentage of voting power to vote for the gauge.

The sum of the weights must be less than or equal to 1 (100%).
Once a vote is made, the vote will continue to be reflected in all epochs until a new vote is made.

## Tally of votes

The votes will be tallied at the start of each epoch.

1. Tally the votes by validators. The validator votes include delegated voting power.
1. Tally the non-validated ballots. If the address is delegated to a validator, its own voting power is deducted from the validator's votes.

### Example

- Addresses

Validator A 1000vRISE
Delegator B 200vRISE (100vRISE delegated to Validator A)

- Pools

Liquidity Pool #1
Liquidity Pool #2

- Votes

A votes Pool #1 & #2 (50% & 50%)
B votes Pool #1 (100%)

#### Result

1. Only A voted
Pool #1's voting power: 550vRISE
Pool #2's voting power: 550vRISE

1. A & B voted
Pool #1's voting power: 700vRISE (500 + 200)
Pool #2's voting power: 650vRISE
