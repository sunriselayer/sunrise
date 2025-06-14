package types

func NewMsgAddYield(admin string, creator string, amount int64) *MsgAddYield {
	return &MsgAddYield{
		Admin:   admin,
		Creator: creator,
		Amount:  amount,
	}
}
