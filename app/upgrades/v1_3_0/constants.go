// Package v1_3_0 is the shutdown upgrade.
// It records application state, sweeps IBC vouchers to the recovery account,
// and leaves the block height to a later governance proposal.
package v1_3_0

// UpgradeName is the governance plan name. The height is chosen when the proposal is submitted.
const UpgradeName = "v1.3.0"

// RecoveryAddress receives every ibc/ voucher except transfer-escrow balances.
const RecoveryAddress = "sunrise1xxgjt7yqkmn63m2d0nrf0vt5uuc2hr6l45xaa9"

// PreUpgradeFileName is the file written under <home>/shutdown before vouchers move.
const PreUpgradeFileName = "pre-upgrade-state.json"

// IBCDenomPrefix is the bank denom prefix of ICS-20 vouchers on this chain.
const IBCDenomPrefix = "ibc/"
