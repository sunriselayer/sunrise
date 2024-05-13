# blobgrant

## Basic architecture

- Users register two addresses
  - `liquidity_provider`: The address possesses the liquidity provider tokens
  - `grantee`: The address to receive DA usage grants
    - One `liquidity_provider` can register one `grantee`
- `x/blobgrant` periodically run ABCI EndBlocker to distribute the fee grant through `x/feegrant`
  - The grant has an expiry date (using standard function of `x/feegrant`)
  - The grant token denom is `blobgrant/gas`.
    - This is non transferrable token
    - `blobgrant/gas` will be converted to the certain amount of protocol revenue in the future.
  - In the antehandler, it is forced to accept `blobgrant/gas` as a fee token only when `x/feegrant` is used in that tx.

## TODO

- Consider how to handle the detection of liquidity provider tokens after the function of "staking liquidity provider tokens for validators"
