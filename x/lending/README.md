# Lending Module

The lending module enables users to lend tokens and receive yield-bearing receipt tokens (riseXXX format), and supports collateralized borrowing against concentrated liquidity positions.

## Features

- **Supply**: Deposit tokens and receive rise tokens (e.g., deposit USDC, receive riseUSDC)
- **Borrow**: Borrow tokens using concentrated liquidity positions as collateral
- **Repay**: Repay borrowed tokens
- **Liquidate**: Liquidate undercollateralized positions

## State

### Markets
- Tracks total supplied and borrowed amounts per token denomination
- Maintains global reward index for yield distribution
- Maps token denoms to their corresponding rise token denoms

### User Positions
- Tracks user's supplied amount per market
- Maintains user's last reward index for yield calculation

### Borrows
- Tracks individual borrow positions with unique IDs
- Links to collateral (liquidity pool ID and position ID)
- Records block height for interest calculation

## Messages

### MsgSupply
Deposit tokens into the lending pool:
- Automatically creates market if it doesn't exist
- Mints rise tokens to the depositor
- Updates market totals and user position

### MsgBorrow
Borrow tokens against collateral:
- Validates collateral value meets LTV requirements
- Checks available liquidity in the market
- Creates borrow record with collateral reference

### MsgRepay
Repay borrowed tokens:
- Supports partial or full repayment
- Updates market totals
- Releases collateral on full repayment

### MsgLiquidate
Liquidate undercollateralized positions:
- Calculates health factor using collateral value
- Transfers liquidation payment from liquidator
- Provides liquidation bonus to liquidator

## Parameters

- **LTV Ratio**: Maximum loan-to-value ratio (default: 80%)
- **Liquidation Threshold**: Health factor below which positions can be liquidated (default: 85%)
- **Base Interest Rate**: Base rate for borrowing (default: 5%)

## Queries

- `Market`: Query specific market by denom
- `Markets`: Query all markets with pagination
- `UserPosition`: Query user's position in specific market
- `UserPositions`: Query all positions for a user
- `Borrow`: Query specific borrow by ID
- `UserBorrows`: Query all borrows for a user

## TODO

- Implement TWAP oracle for accurate price calculations
- Integrate with liquidity pool module for collateral valuation
- Implement interest rate calculations and accrual
- Add yield distribution mechanism for rise token holders
- Implement custom bank module with BeforeSendHook for yield-bearing token transfers