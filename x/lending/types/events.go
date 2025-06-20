package types

// Event types
const (
	EventTypeSupply    = "supply"
	EventTypeBorrow    = "borrow"
	EventTypeRepay     = "repay"
	EventTypeLiquidate = "liquidate"

	AttributeKeySender               = "sender"
	AttributeKeyDenom                = "denom"
	AttributeKeyAmount               = "amount"
	AttributeKeyRiseAmount           = "rise_amount"
	AttributeKeyRiseDenom            = "rise_denom"
	AttributeKeyBorrowId             = "borrow_id"
	AttributeKeyBorrower             = "borrower"
	AttributeKeyCollateralPoolId     = "collateral_pool_id"
	AttributeKeyCollateralPositionId = "collateral_position_id"
	AttributeKeyLiquidator           = "liquidator"
)