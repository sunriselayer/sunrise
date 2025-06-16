package types

func NewMsgClaimCollateralYield(admin string, tokenCreator string, baseYbtCreator string) *MsgClaimCollateralYield {
	return &MsgClaimCollateralYield{
		Admin:          admin,
		TokenCreator:   tokenCreator,
		BaseYbtCreator: baseYbtCreator,
	}
}
