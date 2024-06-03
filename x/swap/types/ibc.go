package types

import (
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"

	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
)

const DefaultRetryCount uint8 = 3

type PacketMetadata struct {
	Swap *SwapMetadata `json:"swap"`
}

type SwapMetadata struct {
	InterfaceProvider string `json:"interface_provider,omitempty"`
	Route             Route  `json:"route,omitempty"`
	ExactAmountIn     *struct {
		MinAmountOut sdkmath.Int `json:"min_amount_out,omitempty"`
	} `json:"exact_amount_in,omitempty"`
	ExactAmountOut *struct {
		AmountOut sdkmath.Int                         `json:"amount_out"`
		Return    *packetforwardtypes.ForwardMetadata `json:"return,omitempty"`
	} `json:"exact_amount_out,omitempty"`
	Forward *packetforwardtypes.ForwardMetadata `json:"forward,omitempty"`
}

func (m *SwapMetadata) Validate() error {
	if err := m.Route.Validate(); err != nil {
		return err
	}

	if m.ExactAmountIn != nil && m.ExactAmountOut != nil {
		return fmt.Errorf("cannot have both exact_amount_in and exact_amount_out")
	}

	if m.ExactAmountIn == nil && m.ExactAmountOut == nil {
		return fmt.Errorf("must have either exact_amount_in or exact_amount_out")
	}

	if m.ExactAmountIn != nil {
		if !m.ExactAmountIn.MinAmountOut.IsPositive() {
			return fmt.Errorf("min_amount_out must be positive")
		}
	}

	if m.ExactAmountOut != nil {
		if m.ExactAmountOut.Return != nil {
			if err := m.ExactAmountOut.Return.Validate(); err != nil {
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
	ForwardAck  []byte      `json:"forward_ack,omitempty"`
	ReturnAck   []byte      `json:"return_ack,omitempty"`
}

func (a SwapAcknowledgement) Acknowledgement() ([]byte, error) {
	bz, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return bz, nil
}
