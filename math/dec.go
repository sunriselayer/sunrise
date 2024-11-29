package math

import (
	sdkmath "cosmossdk.io/math"
)

type Dec struct {
	sdkmath.LegacyDec
}

func NewDec(i int64) Dec {
	return Dec{
		sdkmath.LegacyNewDec(i),
	}
}

func NewDecFromStr(s string) (Dec, error) {
	legacyDec, err := sdkmath.LegacyNewDecFromStr(s)
	if err != nil {
		return Dec{}, err
	}

	return Dec{
		legacyDec,
	}, nil
}

// Marshal implements the gogo proto custom type interface.
func (d Dec) Marshal() ([]byte, error) {
	str := d.LegacyDec.String()
	return []byte(str), nil
}

// MarshalTo implements the gogo proto custom type interface.
func (d *Dec) MarshalTo(data []byte) (n int, err error) {
	bz, err := d.Marshal()
	if err != nil {
		return 0, err
	}

	copy(data, bz)
	return len(bz), nil
}

// Unmarshal implements the gogo proto custom type interface.
func (d *Dec) Unmarshal(data []byte) error {
	if len(data) == 0 {
		d = nil
		return nil
	}

	var err error
	d.LegacyDec, err = sdkmath.LegacyNewDecFromStr(string(data))
	if err != nil {
		return err
	}

	return nil
}

// Size implements the gogo proto custom type interface.
func (d *Dec) Size() int {
	bz, _ := d.Marshal()
	return len(bz)
}

// Override Amino binary serialization by proxying to protobuf.
func (d Dec) MarshalAmino() ([]byte, error)   { return d.Marshal() }
func (d *Dec) UnmarshalAmino(bz []byte) error { return d.Unmarshal(bz) }
