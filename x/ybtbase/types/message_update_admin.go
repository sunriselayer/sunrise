package types

func NewMsgUpdateAdmin(admin string, newAdmin string) *MsgUpdateAdmin {
	return &MsgUpdateAdmin{
		Admin:    admin,
		NewAdmin: newAdmin,
	}
}
