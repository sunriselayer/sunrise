package types

// Event types
const (
	EventTypeSupply    = "supply"
	EventTypeBorrow    = "borrow"
	EventTypeRepay     = "repay"
	EventTypeLiquidate = "liquidate"

	AttributeKeySender     = "sender"
	AttributeKeyDenom      = "denom"
	AttributeKeyAmount     = "amount"
	AttributeKeyRiseAmount = "rise_amount"
	AttributeKeyBorrowId   = "borrow_id"
	AttributeKeyBorrower   = "borrower"
	AttributeKeyLiquidator = "liquidator"
)