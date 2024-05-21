# tokenconverter

This module is for converting `SSR` to `SR` token equivalently if the following rule satisfies

$$
  \text{if} \ \text{CurrentSupplySR} + \text{OutputSR} \le \text{MaxSupplySR}
$$

## Ante Handler

Only for `BlobTx`, this module accept SSR paid for fee.
