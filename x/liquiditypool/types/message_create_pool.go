package types

func NewMsgCreatePool(sender string, denomBase string, denomQuote string) *MsgCreatePool {
	return &MsgCreatePool{
		Sender:     sender,
		DenomBase:  denomBase,
		DenomQuote: denomQuote,
	}
}
