# fee

The module `x/fee` serves the functionalities to burn $RISE tokens used as fees.

- This module makes ante handler only accept tx fee with `fee_denom`.
  - The amount multiplied with `burn_ratio` will be burnt after the tx fee deduction.
- `bypass_denoms` bypass the denom filter and burn.
