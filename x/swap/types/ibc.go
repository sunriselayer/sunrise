package types

import (
	"encoding/json"
	"fmt"
	"strings"

	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	host "github.com/cosmos/ibc-go/v10/modules/core/24-host"
	"github.com/gogo/protobuf/jsonpb"
)

const DefaultRetryCount uint8 = 3

func (m *ForwardMetadata) Validate() error {
	if m.Receiver == "" {
		return fmt.Errorf("failed to validate metadata. receiver cannot be empty")
	}
	if err := host.PortIdentifierValidator(m.Port); err != nil {
		return fmt.Errorf("failed to validate metadata: %w", err)
	}
	if err := host.ChannelIdentifierValidator(m.Channel); err != nil {
		return fmt.Errorf("failed to validate metadata: %w", err)
	}

	return nil
}

func (m *SwapMetadata) Validate() error {
	if err := m.Route.Validate(); err != nil {
		return err
	}
	switch amountStrategy := m.AmountStrategy.(type) {
	case *SwapMetadata_ExactAmountIn:
		if !amountStrategy.ExactAmountIn.MinAmountOut.IsPositive() {
			return fmt.Errorf("min amount out must be positive")
		}

	case *SwapMetadata_ExactAmountOut:
		if amountStrategy.ExactAmountOut.Change != nil {
			if err := amountStrategy.ExactAmountOut.Change.Validate(); err != nil {
				return err
			}
		}
	}

	if m.Forward != nil {
		if err := m.Forward.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type SwapAcknowledgement struct {
	Result      RouteResult `json:"result"`
	IncomingAck []byte      `json:"ibc_ack"`
	ChangeAck   []byte      `json:"change_ack,omitempty"`
	ForwardAck  []byte      `json:"forward_ack,omitempty"`
}

func (a SwapAcknowledgement) Acknowledgement() ([]byte, error) {
	bz, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// GetDenomPrefix returns the receiving denomination prefix
func GetDenomPrefix(portID, channelID string) string {
	return fmt.Sprintf("%s/%s/", portID, channelID)
}

// GetPrefixedDenom returns the denomination with the portID and channelID prefixed
func GetPrefixedDenom(portID, channelID, baseDenom string) string {
	return fmt.Sprintf("%s/%s/%s", portID, channelID, baseDenom)
}

func GetDenomForThisChain(port, channel, counterpartyPort, counterpartyChannel, denom string) string {
	counterpartyPrefix := GetDenomPrefix(counterpartyPort, counterpartyChannel)
	if strings.HasPrefix(denom, counterpartyPrefix) {
		// unwind denom
		unwoundDenom := denom[len(counterpartyPrefix):]
		denomTrace := transfertypes.ExtractDenomFromPath(unwoundDenom)
		if denomTrace.Base == unwoundDenom {
			// denom is now unwound back to native denom
			return unwoundDenom
		}
		// denom is still IBC denom
		return denomTrace.IBCDenom()
	}
	// append port and channel from this chain to denom
	prefixedDenom := GetDenomPrefix(port, channel) + denom
	return transfertypes.ExtractDenomFromPath(prefixedDenom).IBCDenom()
}

func DecodeSwapMetadata(memo string) (*PacketMetadata, error) {
	d := make(map[string]interface{})
	err := json.Unmarshal([]byte(memo), &d)
	if err != nil {
		return nil, err
	}
	if d["swap"] == nil {
		return nil, fmt.Errorf("no swap filed in memo")
	}

	nextString := ""
	swap := d["swap"].(map[string]interface{})
	if swap["forward"] != nil {
		forward := swap["forward"].(map[string]interface{})
		if forward["next"] != nil {
			next := forward["next"]
			delete(forward, "next")
			nextJSON, err := json.Marshal(next)
			if err != nil {
				return nil, err
			}
			nextString = string(nextJSON)

			memoExceptNext, err := json.Marshal(d)
			if err != nil {
				return nil, err
			}
			memo = string(memoExceptNext)
		}
	}

	m := &PacketMetadata{}
	err = jsonpb.Unmarshal(strings.NewReader(memo), m)
	if err != nil {
		return nil, err
	}
	if m.Swap.Forward != nil {
		m.Swap.Forward.Next = nextString
	}
	return m, nil
}
