package types

func NewMsgGrantYieldPermission(admin string, tokenCreator string, target string) *MsgGrantYieldPermission {
	return &MsgGrantYieldPermission{
		Admin:        admin,
		TokenCreator: tokenCreator,
		Target:       target,
	}
}
