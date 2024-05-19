# tokenconverter

This module is for converting `SRG` to `SR` token.

$$
  \text{OutputSR} = \text{InputSSR} \times \min\left(1, \frac{\text{CurrentSupplySR}}{\text{CurrentSupplySSR}} \right) \ \text{if} \ \text{CurrentSupplySR} + \text{OutputSR} \le \text{MaxSupplySR}
$$

## Ante Handler

Only for `BlobTx`, this module accept SSR paid for fee.
