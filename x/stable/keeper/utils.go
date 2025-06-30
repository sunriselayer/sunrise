// Package keeper provides the core logic for the stable module.
// This file contains utility functions used across the keeper's implementation.
package keeper

// abs returns the absolute value of an integer.
// It's a helper function used in calculations involving exponents.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
