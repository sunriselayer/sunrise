package types

func NewMsgBurn(admin string, creator string, amount int64) *MsgBurn {
	return &MsgBurn{
		Admin:   admin,
		Creator: creator,
		Amount:  amount,
	}
}
