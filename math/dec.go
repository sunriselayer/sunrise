package math

import (
	sdkmath "cosmossdk.io/math"
)

type LegacyDec struct {
	sdkmath.LegacyDec
}

func LegacyZeroDec() LegacyDec     { return LegacyDec{sdkmath.LegacyZeroDec()} }
func LegacyOneDec() LegacyDec      { return LegacyDec{sdkmath.LegacyOneDec()} }
func LegacySmallestDec() LegacyDec { return LegacyDec{sdkmath.LegacySmallestDec()} }

func LegacyNewDec(i int64) LegacyDec {
	return LegacyDec{
		sdkmath.LegacyNewDec(i),
	}
}

func LegacyNewDecWithPrec(i, prec int64) LegacyDec {
	return LegacyDec{
		sdkmath.LegacyNewDecWithPrec(i, prec),
	}
}

func LegacyNewDecFromStr(s string) (LegacyDec, error) {
	legacyDec, err := sdkmath.LegacyNewDecFromStr(s)
	if err != nil {
		return LegacyDec{}, err
	}

	return LegacyDec{
		legacyDec,
	}, nil
}

// Marshal implements the gogo proto custom type interface.
func (d LegacyDec) Marshal() ([]byte, error) {
	str := d.LegacyDec.String()
	return []byte(str), nil
}

// MarshalTo implements the gogo proto custom type interface.
func (d *LegacyDec) MarshalTo(data []byte) (n int, err error) {
	bz, err := d.Marshal()
	if err != nil {
		return 0, err
	}

	copy(data, bz)
	return len(bz), nil
}

// Unmarshal implements the gogo proto custom type interface.
func (d *LegacyDec) Unmarshal(data []byte) error {
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
func (d *LegacyDec) Size() int {
	bz, _ := d.Marshal()
	return len(bz)
}

// Override Amino binary serialization by proxying to protobuf.
func (d LegacyDec) MarshalAmino() ([]byte, error)   { return d.Marshal() }
func (d *LegacyDec) UnmarshalAmino(bz []byte) error { return d.Unmarshal(bz) }
