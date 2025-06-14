package types

func NewMsgRevokeYieldPermission(admin string, creator string, target string) *MsgRevokeYieldPermission {
	return &MsgRevokeYieldPermission{
		Admin:   admin,
		Creator: creator,
		Target:  target,
	}
}
