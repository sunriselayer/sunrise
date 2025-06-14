package types

func NewMsgCreate(creator string, admin string, permissioned bool) *MsgCreate {
	return &MsgCreate{
		Creator:      creator,
		Admin:        admin,
		Permissioned: permissioned,
	}
}
