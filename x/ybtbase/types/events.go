package types

// Event types
const (
	EventTypeCreateToken      = "create_token"
	EventTypeMint             = "mint"
	EventTypeBurn             = "burn"
	EventTypeAddYield         = "add_yield"
	EventTypeClaimYield       = "claim_yield"
	EventTypeGrantPermission  = "grant_yield_permission"
	EventTypeRevokePermission = "revoke_yield_permission"
	EventTypeUpdateAdmin      = "update_admin"

	// Attribute keys
	AttributeKeyCreator      = "creator"
	AttributeKeyAdmin        = "admin"
	AttributeKeyNewAdmin     = "new_admin"
	AttributeKeyPermissioned = "permissioned"
	AttributeKeyAmount       = "amount"
	AttributeKeyTarget       = "target"
	AttributeKeyDenom        = "denom"
	AttributeKeyClaimer      = "claimer"
	AttributeKeyYieldAmount  = "yield_amount"
)
