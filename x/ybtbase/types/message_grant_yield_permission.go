package types

func NewMsgGrantYieldPermission(admin string, creator string, target string) *MsgGrantYieldPermission {
	return &MsgGrantYieldPermission{
		Admin:   admin,
		Creator: creator,
		Target:  target,
	}
}
