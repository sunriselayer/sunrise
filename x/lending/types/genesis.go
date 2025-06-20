package types

import (
	"fmt"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		Markets:       []Market{},
		UserPositions: []UserPosition{},
		Borrows:       []Borrow{},
		BorrowCount:   0,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Check for duplicate markets
	marketMap := make(map[string]bool)
	for _, market := range gs.Markets {
		if _, exists := marketMap[market.Denom]; exists {
			return fmt.Errorf("duplicate market denom: %s", market.Denom)
		}
		marketMap[market.Denom] = true

		// Validate market fields
		if market.Denom == "" {
			return fmt.Errorf("market denom cannot be empty")
		}
		if market.RiseDenom == "" {
			return fmt.Errorf("rise denom cannot be empty for market %s", market.Denom)
		}
		if market.TotalSupplied.IsNegative() {
			return fmt.Errorf("total supplied cannot be negative for market %s", market.Denom)
		}
		if market.TotalBorrowed.IsNegative() {
			return fmt.Errorf("total borrowed cannot be negative for market %s", market.Denom)
		}
		if market.GlobalRewardIndex.IsNegative() {
			return fmt.Errorf("global reward index cannot be negative for market %s", market.Denom)
		}
	}

	// Check for duplicate user positions
	type positionKey struct {
		user  string
		denom string
	}
	positionMap := make(map[positionKey]bool)
	for _, position := range gs.UserPositions {
		key := positionKey{user: position.UserAddress, denom: position.Denom}
		if _, exists := positionMap[key]; exists {
			return fmt.Errorf("duplicate user position: %s, %s", position.UserAddress, position.Denom)
		}
		positionMap[key] = true

		// Validate position fields
		if position.UserAddress == "" {
			return fmt.Errorf("user address cannot be empty")
		}
		if position.Denom == "" {
			return fmt.Errorf("position denom cannot be empty")
		}
		if position.Amount.IsNegative() {
			return fmt.Errorf("position amount cannot be negative")
		}
		if position.LastRewardIndex.IsNegative() {
			return fmt.Errorf("last reward index cannot be negative")
		}

		// Check that market exists for this position
		if _, exists := marketMap[position.Denom]; !exists {
			return fmt.Errorf("position references non-existent market: %s", position.Denom)
		}
	}

	// Check for duplicate borrows and validate
	borrowMap := make(map[uint64]bool)
	maxBorrowId := uint64(0)
	for _, borrow := range gs.Borrows {
		if _, exists := borrowMap[borrow.Id]; exists {
			return fmt.Errorf("duplicate borrow id: %d", borrow.Id)
		}
		borrowMap[borrow.Id] = true

		// Track max borrow ID
		if borrow.Id > maxBorrowId {
			maxBorrowId = borrow.Id
		}

		// Validate borrow fields
		if borrow.Borrower == "" {
			return fmt.Errorf("borrower address cannot be empty for borrow %d", borrow.Id)
		}
		if !borrow.Amount.IsValid() || borrow.Amount.IsZero() {
			return fmt.Errorf("invalid borrow amount for borrow %d", borrow.Id)
		}
		if borrow.CollateralPoolId == 0 {
			return fmt.Errorf("collateral pool id cannot be zero for borrow %d", borrow.Id)
		}
		if borrow.CollateralPositionId == 0 {
			return fmt.Errorf("collateral position id cannot be zero for borrow %d", borrow.Id)
		}
		if borrow.BlockHeight <= 0 {
			return fmt.Errorf("block height must be positive for borrow %d", borrow.Id)
		}

		// Check that market exists for this borrow
		if _, exists := marketMap[borrow.Amount.Denom]; !exists {
			return fmt.Errorf("borrow references non-existent market: %s", borrow.Amount.Denom)
		}
	}

	// Validate borrow count
	if gs.BorrowCount < maxBorrowId {
		return fmt.Errorf("borrow count %d is less than max borrow id %d", gs.BorrowCount, maxBorrowId)
	}

	return nil
}
