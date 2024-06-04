# swap

## MsgSwapExactAmountIn

## MsgSwapExactAmountOut

## ICS20 Middleware

Swap functions also can be executed by ICS20 token transfer packet automatically.

### Metadata

JSON string of marshalled `PacketMetadata` should be inserted in the `memo` field of ICS20 transfer packet.

```typescript
type PacketMetadata = {
  [namespace: string]: unknown;
  swap?: SwapMetadata;
};

type SwapMetadata = {
  interface_provider: string;
  route: Route;

  forward?: ForwardMetadata;
} & (
  | {
      exact_amount_in: {
        min_amount_out: string;
      };
    }
  | {
      exact_amount_out: {
        amount_out: string;
        change?: ForwardMetadata;
      };
    }
);

type ForwardMetadata = {
  receiver: string;
  port: string;
  channel: string;
  timeout: string;
  retries: number;
  next?: PacketMetadata;
};
```

`ForwardMetadata` is quoted from [Packet Forward Middleware](https://github.com/cosmos/ibc-apps/tree/main/middleware/packet-forward-middleware).

### Sequence diagrams

#### Neither Return nor Forward

```mermaid
sequenceDiagram
    autonumber
    Chain A ->> Sunrise: Transfer token X
    Sunrise --> Sunrise: recv_packet
    Sunrise ->> Sunrise: Swap token X to token Y
    Sunrise ->> Chain A: ack
```

#### Forward

```mermaid
sequenceDiagram
    autonumber
    Chain A ->> Sunrise: Transfer token X
    Sunrise --> Sunrise: recv_packet
    Sunrise ->> Sunrise: Swap token X to token Y

    Sunrise ->> Chain B: Forward token Y
    Chain B --> Chain B: recv_packet
    Chain B ->> Sunrise: ack
    Sunrise ->> Chain A: ack
```

#### Change and Forward

If the exact output amount is designated for the swap, the remainder input amount will occur.
There is a function to automatically refund the remainder input amount.

```mermaid
sequenceDiagram
    autonumber
    Chain A ->> Sunrise: Transfer token X
    Sunrise --> Sunrise: recv_packet
    Sunrise ->> Sunrise: Swap token X to token Y

    Sunrise ->> Chain A: Change token X
    Sunrise ->> Chain B: Forward token Y
    Chain A --> Chain A: recv_packet
    Chain B --> Chain B: recv_packet
    Chain A ->> Sunrise: ack
    Chain B ->> Sunrise: ack
    Sunrise ->> Chain A: ack
```

### Receiver address

After the swapping has been executed, the acknowledgement of "Transfer token X" will be always success even if the next return / forward packet failed. The swapped funds are preserved in the balance of the receiver address.
