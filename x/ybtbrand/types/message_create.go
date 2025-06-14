package types

func NewMsgCreate(creator string, admin string) *MsgCreate {
	return &MsgCreate{
		Creator: creator,
		Admin:   admin,
	}
}
