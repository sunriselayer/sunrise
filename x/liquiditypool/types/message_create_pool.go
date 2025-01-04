package types

func NewMsgCreatePool(authority string) *MsgCreatePool {
	return &MsgCreatePool{
		Authority: authority,
	}
}
