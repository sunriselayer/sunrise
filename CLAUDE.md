# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sunrise is a Cosmos SDK-based blockchain focused on data availability and DeFi functionality. It uses Cosmos SDK v0.53.2, IBC-Go v10.2.0, and includes custom modules for data availability (with ZK proofs), liquidity pools, swaps, and incentive mechanisms.

## Key Commands

### Build and Development

```bash
# Install the sunrised binary
make install

# Run development server with hot reload (recommended)
ignite chain serve

# Build without installing
go build -v ./...

# Generate protobuf code after modifying .proto files
make proto-gen
```

### Testing

```bash
# Run all tests (includes govet and vulnerability check)
make test

# Run unit tests only
make test-unit

# Run tests with race detection
make test-race

# Run tests with coverage
make test-cover

# Run benchmarks
make bench
```

### Code Quality

```bash
# Run linter
make lint

# Run linter with auto-fix
make lint-fix

# Check for vulnerabilities
make govulncheck
```

## Architecture

### Module Structure

The blockchain uses a modular architecture with custom modules in `/x/`:

- **bank**: Custom bank module with BeforeSendHook for yield-bearing token support
- **da**: Data availability module implementing erasure coding and zero-knowledge proofs for data availability attestations
- **fee**: Custom fee management beyond standard Cosmos SDK fees
- **lending**: Lending protocol with yield-bearing receipt tokens (planned)
- **liquidityincentive**: Incentive mechanism with epochs, gauges, and bribe functionality
- **liquiditypool**: AMM-style liquidity pool management
- **lockup**: Token lockup/vesting functionality
- **rfq**: Request-For-Quote module for cross-chain swaps (planned)
- **shareclass**: Share classification system for different token types
- **swap**: DEX functionality with routing capabilities
- **tokenconverter**: Cross-token conversion mechanisms

### Key Integration Points

- **IBC**: Full IBC support with WASM light client capability for cross-chain communication (uses WasmVM for light clients only)
- **Ante Handlers**: Custom transaction processing in `app/ante/` for fee and swap handling
- **Vote Extensions**: Enabled from genesis for advanced consensus features

### Development Workflow

1. Protocol changes start in `/proto/sunrise/[module]/` files
2. Run `make proto-gen` to generate Go code
3. Implement keeper methods in `/x/[module]/keeper/`
4. Add message handlers in `/x/[module]/keeper/msg_server.go`
5. Write tests alongside implementation
6. Integration tests go in `/tests/e2e/` or `/tests/interchain/`

### Testing Patterns

- Unit tests use standard Go testing with mocks generated via `mockgen`
- Keeper tests typically use a test keeper setup (see `x/*/keeper/keeper_test.go`)
- E2E tests use the interchain test framework for multi-chain scenarios
- Always run `make test` before committing to catch linting and vulnerability issues

## Important Configurations

### Genesis Parameters (Development)

- Voting period: 20s (for fast testing)
- DA module periods: 30-60s
- Vote extensions: Enabled from height 1

### Token Denominations

- Fee token: `urise`
- Governance token: `uvrise`

### Development Accounts

The development chain (via `ignite chain serve`) includes pre-funded test accounts: `val`, `faucet`, `user1-4`.

## TypeScript Client

Auto-generated TypeScript client is in `/ts-client/`. Regenerate with `ignite chain serve` or after protocol buffer changes.

## Module Implementation Status

### x/lending

A lending module that enables yield-bearing token creation through lending:

- **Purpose**: Allow depositors to lend tokens and receive yield-bearing receipt tokens (riseXXX format)
- **Status**: ✅ Core functionality implemented with TDD approach

**Implemented Features**:

1. **Message Handlers** (All tested):
   - `MsgSupply`: Deposit tokens and receive riseXXX tokens (auto-creates markets)
   - `MsgBorrow`: Borrow using concentrated liquidity positions as collateral
   - `MsgRepay`: Repay borrowed tokens (partial/full)
   - `MsgLiquidate`: Liquidate undercollateralized positions
   - `MsgUpdateParams`: Update module parameters (gov only)

2. **State Management**:
   - Markets: Track total supplied/borrowed, global reward index, rise denom mapping
   - UserPositions: Track user's supplied amount and last reward index
   - Borrows: Individual borrow records with collateral references
   - BorrowId: Auto-incrementing sequence for unique borrow IDs

3. **Query Endpoints**:
   - Market/Markets: Query lending markets with pagination
   - UserPosition/UserPositions: Query user's lending positions
   - Borrow/UserBorrows: Query borrow information
   - Params: Query module parameters (LTV, liquidation threshold, interest rate)

4. **Genesis Import/Export**:
   - Full state import/export functionality
   - Comprehensive validation of genesis state
   - Support for all state objects

**TODO - Core Features**:

- [ ] **TWAP Oracle Integration** (currently mocked)
  - Integrate with oracle module for accurate price feeds
  - Support multiple price sources for reliability
- [ ] **Liquidity Pool Integration** 
  - Connect with liquiditypool module for collateral valuation
  - Implement position value calculations
- [ ] **Interest Rate Model**
  - Implement dynamic interest rates based on utilization
  - Add interest accrual on each block
  - Update global reward index periodically
- [ ] **Yield Distribution for rise Tokens**
  - Implement yield accounting in custom bank module
  - Add BeforeSendHook for transfer handling
  - Support negative userLastRewardIndex for transfers
- [ ] **Liquidation Mechanics**
  - Calculate liquidation bonus for liquidators
  - Implement collateral auction mechanism
  - Add partial liquidation support

**TODO - Additional Features**:

- [ ] CLI commands for all operations
- [ ] Emergency pause functionality
- [ ] Protocol fee collection
- [ ] IBC token integration tests
- [ ] Monitoring and metrics hooks

### x/rfq

Request-For-Quote (RFQ) module for efficient cross-chain swaps:

- **Purpose**: Enable capital-efficient market making without requiring fillers to hold large inventories
- **Status**: 🔄 Not yet implemented

**Planned Architecture**:
- RFQ-based price discovery mechanism
- Integration with lending module for collateralized borrowing
- Cross-chain message passing via IBC
- Reputation system for fillers

**TODO - Implementation**:
- [ ] Design proto messages for RFQ lifecycle
- [ ] Implement quote submission and matching
- [ ] Add filler registration system
- [ ] Integrate with lending for capital efficiency
- [ ] Create cross-chain swap execution flow
- [ ] Add slashing for failed swaps
- [ ] Implement keeper and message handlers
- [ ] Add comprehensive query endpoints
- [ ] Create CLI commands
- [ ] Write unit and integration tests

### Integration Architecture

**Yield-Bearing Token Flow**:

1. User supplies USDC to lending module
2. Lending module mints riseUSDC tokens
3. Custom bank module tracks yield accrual on transfers
4. Users can claim accumulated yield

**RFQ-Lending Integration**:

1. Filler deposits concentrated liquidity position as collateral
2. Borrows tokens needed for RFQ fill
3. Executes cross-chain swap
4. Repays loan with tokens from other chain
5. Collateral released after successful repayment

**Risk Management**:

- Segregated risk between lenders and fillers
- Over-collateralization requirements
- Liquidation mechanisms for bad debt
- Time-bound repayment for RFQ borrows
