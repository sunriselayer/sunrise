package types

func NewMsgCreatePool(creator string) *MsgCreatePool {
	return &MsgCreatePool{
		Creator: creator,
	}
}
