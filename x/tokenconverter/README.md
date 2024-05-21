# tokenconverter

This module is for converting `SSR` to `SR` token.

$$
  \text{OutputSR} = \text{InputSSR} \times \min\left(1, \frac{\text{MaxSupplySR}}{\text{CurrentSupplySSR}} \right) \ \text{if} \ \text{CurrentSupplySR} + \text{OutputSR} \le \text{MaxSupplySR}
$$

## Ante Handler

Only for `BlobTx`, this module accept SSR paid for fee.
