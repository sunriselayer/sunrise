# DA

## Specification for Zero-Knowledge Proof

### Terms and Notation

- $n$: Total number of shards
- $t$: Threshold (minimum number of shards required to prove the possession)
- $s_i$: The $i$-th shard
- $H$: The hash function

### Overview

This system verifies the possession of data shard hash `H(s_i)` without exposing `H(s_i)`.

### Zero-Knowledge Proof System

#### 4.1 Public Inputs

- $\{H^2(s_i)\}_{i=1}^n$
- $t$: Threshold

#### 4.2 Private Inputs

- $\{H(s_i)\}_{i=1}^t$

#### 4.3 ZKP Circuit Constraints

Hash Verification:
For each $i$, verify $H^2(s_i) = H^2(s_i)_{public}$