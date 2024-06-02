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
	Routes       []Route                             `json:"routes,omitempty"`
	MinAmountOut sdkmath.Int                         `json:"min_amount_out,omitempty"`
	Forward      *packetforwardtypes.ForwardMetadata `json:"forward,omitempty"`
}

func (m *SwapMetadata) Validate() error {
	if len(m.Routes) == 0 {
		return fmt.Errorf("failed to validate metadata. routes cannot be empty")
	}

	if m.Forward != nil {
		err := m.Forward.Validate()
		if err != nil {
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
