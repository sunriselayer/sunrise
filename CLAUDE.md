# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sunrise is a Cosmos SDK-based blockchain focused on data availability and DeFi functionality. The project uses Go 1.23.4+ and the Ignite CLI for development.

## Development Commands

### Building and Installation
- `make install` - Build and install the `sunrised` binary
- `make all` - Alias for install

### Testing
- `make test` - Run complete test suite (includes govet, govulncheck, and unit tests)
- `make test-unit` - Run unit tests only with 30m timeout
- `make test-race` - Run tests with race condition detection
- `make test-cover` - Generate coverage report
- `make bench` - Run benchmarks
- `go test ./x/[module]/...` - Run tests for specific module

### Code Quality
- `make lint` - Run golangci-lint with 15m timeout
- `make lint-fix` - Run linter and auto-fix issues
- `make govet` - Run go vet
- `make govulncheck` - Run vulnerability check

### Protobuf Generation
- `make proto-gen` - Generate protobuf files using Ignite CLI
- `ignite generate proto-go --yes` - Alternative protobuf generation

### Development Server
- `ignite chain serve` - Start development blockchain (installs deps, builds, initializes, starts)

## Architecture

### Core Modules (x/ directory)
The blockchain includes several custom modules:

- **da** - Data Availability layer with zero-knowledge proofs for data possession verification
- **liquiditypool** - Concentrated liquidity AMM with tick-based pricing (similar to Uniswap V3)  
- **swap** - Token swapping with complex routing (series/parallel routes) and IBC middleware integration
- **liquidityincentive** - Bribe-based liquidity incentives with epoch-based voting
- **shareclass** - Validator share management and delegation rewards
- **lockup** - Token lockup functionality with vesting
- **fee** - Custom fee handling and burning mechanisms
- **tokenconverter** - Token conversion between different formats

### Key Design Patterns

#### Module Structure
Each custom module follows standard Cosmos SDK patterns:
- `keeper/` - Core business logic and state management
- `types/` - Message types, genesis, params, and protobuf generated code
- `module/` - Module interface implementation (autocli, genesis, simulation)
- `testutil/` - Mock interfaces for testing (when needed)

#### Testing Approach
- Uses `go.uber.org/mock/gomock` for mocking
- Test files follow `*_test.go` naming in `keeper_test` package
- Fixture pattern for test setup with mocks and context
- 30-minute timeout for all tests

#### Data Availability (DA) Module
- Implements erasure coding with zero-knowledge proofs
- Uses challenge-response mechanism for data verification  
- Supports IPFS and Arweave for data storage
- Validator voting and slashing for data availability consensus

#### Liquidity Pool Architecture
- Concentrated liquidity with configurable tick parameters
- Price calculation: `price(tick) = price_ratio^(tick - base_offset)`
- Position-based liquidity provision
- Integration with swap module for routing

#### IBC Integration
- Swap module acts as IBC middleware for automatic token swapping
- Supports packet forwarding with metadata-driven routing
- JSON metadata in ICS20 transfer memo field

### Configuration
- `config.yml` - Ignite development configuration with validators and faucet setup
- Uses `urise` and `uvrise` denominations
- Default account prefix: `sunrise`

## Dependencies and Tools

### Key Dependencies
- Cosmos SDK v0.53.2
- CometBFT v0.38.17
- IBC-Go v10.2.0
- gnark v0.12.0 (for zero-knowledge proofs)

### Development Tools
- Ignite CLI (nightly build for Cosmos SDK v0.52+ compatibility)
- golangci-lint v1.61.0
- buf for protobuf management
- gomock for test mocking