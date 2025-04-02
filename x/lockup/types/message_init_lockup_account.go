package types

func NewMsgInitLockupAccount(sender string) *MsgInitLockupAccount {
	return &MsgInitLockupAccount{
		Sender: sender,
	}
}
