package types

// Event types
const (
	EventTypeCreate               = "create_ybtbrand"
	EventTypeMint                 = "mint_ybtbrand"
	EventTypeBurn                 = "burn_ybtbrand"
	EventTypeAddYields            = "add_yields"
	EventTypeClaimYields          = "claim_yields"
	EventTypeUpdateAdmin          = "update_admin"
	EventTypeClaimCollateralYield = "claim_collateral_yield"

	AttributeKeyCreator        = "creator"
	AttributeKeyAdmin          = "admin"
	AttributeKeyNewAdmin       = "new_admin"
	AttributeKeyBaseYbtCreator = "base_ybt_creator"
	AttributeKeyDenom          = "denom"
	AttributeKeyAmount         = "amount"
	AttributeKeyYieldAmount    = "yield_amount"
	AttributeKeyYieldDenom     = "yield_denom"
	AttributeKeyClaimer        = "claimer"
)
