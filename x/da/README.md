# DA

## Overview

- Offloading erasure coding to off-chain workers
- Offloading blob delivery from mempool
  - Network bandwidth consumption is reduced to constant order from linear with blob size
- Deterministic shard allocation for each validator
- KZG-commitment based blazing fast data availability
  - Minimum 2 blocks for confirmation
    - One block for `MsgDeclareBlob`
    - One block for `MsgBundleCommitments`
