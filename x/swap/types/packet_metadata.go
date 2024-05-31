package types

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
)

type PacketMetadata struct {
	Swap *SwapMetadata `json:"swap"`
}

type SwapMetadata struct {
	Routes       []Route  `json:"routes,omitempty"`
	MinAmountOut math.Int `json:"min_amount_out,omitempty"`
	Return       *struct {
		Port     string        `json:"port,omitempty"`
		Channel  string        `json:"channel,omitempty"`
		Timeout  time.Duration `json:"timeout,omitempty"`
		Receiver string        `json:"receiver,omitempty"`
		Memo     string        `json:"memo,omitempty"`
	} `json:"return,omitempty"`
}

func (m *SwapMetadata) Validate() error {
	if len(m.Routes) == 0 {
		return fmt.Errorf("failed to validate metadata. routes cannot be empty")
	}

	if m.Return != nil {
		if err := host.PortIdentifierValidator(m.Return.Port); err != nil {
			return fmt.Errorf("failed to validate metadata: %w", err)
		}
		if err := host.ChannelIdentifierValidator(m.Return.Channel); err != nil {
			return fmt.Errorf("failed to validate metadata: %w", err)
		}
		if m.Return.Receiver == "" {
			return fmt.Errorf("failed to validate metadata. receiver cannot be empty")
		}
	}

	return nil
}
