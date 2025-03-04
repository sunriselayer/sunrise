package types

func NewMsgCreatePool(authority string, denomBase string, denomQuote string) *MsgCreatePool {
	return &MsgCreatePool{
		Authority:  authority,
		DenomBase:  denomBase,
		DenomQuote: denomQuote,
	}
}
