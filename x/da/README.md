# DA

## Specification for Zero-Knowledge Proof of Shamir's Secret Sharing Reconstruction

### 1. Overview

This system verifies the reconstruction process of a secret distributed using Shamir's Secret Sharing scheme, utilizing Zero-Knowledge Proofs (ZKP). The system proves that the shares used for reconstruction are valid and that the number of shares used meets or exceeds the threshold.

### 2. Terms and Notation

- $S$: The secret value
- $n$: Total number of shares
- $t$: Threshold (minimum number of shares required to reconstruct the secret)
- $h_i$: The $i$-th share
- $x_i$: The x-coordinate corresponding to the $i$-th share
- $H(x)$: Poseidon hash function

### 3. Shamir's Secret Sharing

#### 3.1 Coefficient Determination

Given a set of shares and their corresponding x-coordinates, we determine the coefficients of the polynomial. This is essentially solving a system of linear equations.

We are given $t$ shares $\{h_i\}_{i=1}^t$ and their corresponding x-coordinates $\{x_i\}_{i=1}^t$.
The polynomial $f(x)$ of degree $t-1$ is defined as:
$f(x) = a_0 + a_1x + a_2x^2 + ... + a_{t-1}x^{t-1} \mod p$
where $a_0 = S$ (the secret), and $a_1, ..., a_{t-1}$ are the coefficients we need to determine.
For each share, we have the equation:
$h_i = f(x_i) = a_0 + a_1x_i + a_2x_i^2 + ... + a_{t-1}x_i^{t-1} \mod p$, for $i = 1, ..., t$
This forms a system of $t$ linear equations with $t$ unknowns $(a_0, a_1, ..., a_{t-1})$. We can solve this system using methods such as Gaussian elimination or matrix inversion.
The solution to this system gives us the coefficients of the polynomial, including $a_0$ which is the secret $S$.

#### 3.2 Secret Reconstruction

Reconstruct the secret using Lagrange interpolation:
$S = \sum_{i=1}^t h_i \cdot \prod_{j \neq i} \frac{x_j}{x_j - x_i} \mod p$

### 4. Zero-Knowledge Proof System

#### 4.1 Public Inputs

- ${H(h_i)}_{i=1}^n$: Poseidon hash values of each share
- $t$: Threshold

#### 4.2 Private Inputs

- ${h_i}_{i=1}^t$: Shares used
- ${x_i}_{i=1}^t$: Corresponding x-coordinates
- $S$: Reconstructed secret

#### 4.3 ZKP Circuit Constraints

Hash Verification:
For each $i$, verify $H(h_i) = H(h_i)_{public}$
Secret Reconstruction:
Compute and verify $S = \sum_{i=1}^t h_i \cdot \prod_{j \neq i} \frac{x_j}{x_j - x_i}$
Threshold Verification:
Verify that the number of shares used $\geq t$

### 5. Security Considerations

The Poseidon hash function is a cryptographic hash function optimized for ZKP systems.
The secret $S$ and each share $h_i$ are operated on in a finite field of appropriate size.
The choice of x-coordinates $x_i$ can affect security and should be done carefully.

### 6. Implementation Notes

The ZKP circuit is implemented using the gnark library.
Large integer operations are handled using the math/big package.
Cryptographically secure random number generators are used for generating randomness.

### 7. Performance Considerations

The complexity of the ZKP circuit increases proportionally with the number of shares $t$ used.
Poseidon hash computations may be one of the most computationally expensive operations in the circuit.

This specification outlines the implementation and key mathematical concepts. Actual implementation may require more detailed specifications, particularly regarding error handling, input validation, and specific parameter choices.
