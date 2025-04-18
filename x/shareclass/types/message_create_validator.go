package types

func NewMsgCreateValidator(validatorAddress string) *MsgCreateValidator {
	return &MsgCreateValidator{
		ValidatorAddress: validatorAddress,
	}
}
