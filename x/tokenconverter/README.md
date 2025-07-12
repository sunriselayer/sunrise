# Token Converter Module

The `x/tokenconverter` module provides a mechanism for converting between transferable and non-transferable tokens.

## Core Concepts

The token converter module allows for the seamless conversion of tokens between two states: transferable and non-transferable. This is useful for scenarios where you want to restrict the transfer of certain tokens, such as those that are locked or vested.

### Conversion

The module provides two conversion functions:

- **Convert**: Converts non-transferable tokens to transferable tokens.
- **ConvertReverse**: Converts transferable tokens to non-transferable tokens.

Both conversions are 1:1, meaning that one non-transferable token is converted to one transferable token, and vice versa.

## Messages

### MsgConvert

Converts non-transferable tokens to transferable tokens.

```protobuf
message MsgConvert {
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  cosmossdk.io.math.Int amount = 2;
}
```

## Parameters

The `x/tokenconverter` module has the following parameters:

- `non_transferable_denom`: The denomination of the non-transferable token.
- `transferable_denom`: The denomination of the transferable token.
