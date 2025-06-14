package types

func NewMsgMint(admin string, creator string, amount int64) *MsgMint {
	return &MsgMint{
		Admin:   admin,
		Creator: creator,
		Amount:  amount,
	}
}
