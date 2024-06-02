package types

import (
	"encoding/json"

	sdkmath "cosmossdk.io/math"

	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
)

const DefaultRetryCount uint8 = 3

type PacketMetadata struct {
	Swap *SwapMetadata `json:"swap"`
}

type SwapMetadata struct {
	Route        Route                               `json:"route,omitempty"`
	MinAmountOut sdkmath.Int                         `json:"min_amount_out,omitempty"`
	Forward      *packetforwardtypes.ForwardMetadata `json:"forward,omitempty"`
}

func (m *SwapMetadata) Validate() error {
	if err := m.Route.Validate(); err != nil {
		return err
	}

	if m.Forward != nil {
		if err := m.Forward.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type SwapAcknowledgement struct {
	AmountOut   sdkmath.Int `json:"amount_out"`
	IncomingAck []byte      `json:"ibc_ack"`
	ForwardAck  []byte      `json:"forward_ack,omitempty"`
}

func (a SwapAcknowledgement) Acknowledgement() ([]byte, error) {
	bz, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return bz, nil
}
