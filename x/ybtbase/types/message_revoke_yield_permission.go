package types

func NewMsgRevokeYieldPermission(admin string, tokenCreator string, target string) *MsgRevokeYieldPermission {
	return &MsgRevokeYieldPermission{
		Admin:        admin,
		TokenCreator: tokenCreator,
		Target:       target,
	}
}
