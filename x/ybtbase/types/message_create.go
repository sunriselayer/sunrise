package types

func NewMsgCreate(creator string, admin string, permissionMode PermissionMode) *MsgCreate {
	return &MsgCreate{
		Creator:        creator,
		Admin:          admin,
		PermissionMode: permissionMode,
	}
}
