# bank

This module overrides the default bank module to add a custom hook for the `Send` function.

## BeforeSendHook

The `Send` function is overridden to call the `BeforeSend` hook.
If the hook returns an error, the transaction will fail.

```go
type BeforeSendHook func(ctx context.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) error
```
