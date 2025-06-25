# Gemini Project Notes: sunrise

This document contains notes and context about the `sunrise` project, learned through our interactions. It is based on information from the `Makefile` and `CLAUDE.md`.

## Project Overview

Sunrise is a Cosmos SDK-based blockchain focused on data availability and advanced DeFi functionality. It's built with Go and utilizes the Ignite CLI for development workflows.

**Key Technologies:**
*   Go 1.23.4+
*   Cosmos SDK v0.53.2
*   CometBFT v0.38.17
*   IBC-Go v10.2.0
*   gnark v0.12.0 (for Zero-Knowledge Proofs)

## Key Commands

Based on the `Makefile`, here are the primary commands for working with this project:

*   **Start Dev Chain:** `ignite chain serve`
    *   The primary command for local development. It handles installation, initialization, and starts the blockchain.
*   **Install/Build:** `make install`
    *   Builds and installs the main application binary `sunrised`.
*   **Run Unit Tests:** `make test-unit`
    *   Executes all unit tests in the project.
*   **Run Comprehensive Tests:** `make test`
    *   A convenient target that runs `go vet`, `govulncheck`, and `make test-unit` sequentially. This should be the standard check before committing code.
*   **Run Linter:** `make lint`
    *   Runs `golangci-lint` to check for code style and quality issues.
*   **Run Vulnerability Check:** `make govulncheck`
    *   Scans the project's dependencies for known vulnerabilities.
*   **Generate Protobuf Files:** `make proto-gen`
    *   Uses `ignite` to generate Go code from the `.proto` files.

## Architecture & Conventions

### Core Modules (`x/` directory)
*   **`da`**: Data Availability layer using erasure coding and ZKPs (gnark) for data possession verification.
*   **`liquiditypool`**: Implements a concentrated liquidity AMM, similar to Uniswap V3, with tick-based pricing.
*   **`swap`**: Handles token swapping with complex routing (series/parallel) and acts as IBC middleware.
*   **`liquidityincentive`**: Manages bribe-based liquidity incentives through epoch-based voting.
*   **`shareclass`**: Manages validator shares and delegation rewards.
*   **`lockup`**: Provides token lockup functionality, including vesting schedules.
*   **`fee`**: Implements custom fee handling and burning mechanisms.
*   **`tokenconverter`**: Handles fixed-rate or other non-AMM token conversions.

### Design Patterns
*   **Module Structure**: Standard Cosmos SDK layout (`keeper`, `types`, `module`).
*   **Testing**: Uses `go.uber.org/mock/gomock` for mocking, with a fixture pattern for setup. Tests have a 30-minute timeout.
*   **IBC Integration**: The `swap` module functions as IBC middleware, enabling automatic swaps for incoming transfers via metadata in the ICS20 memo field.
*   **Configuration**: `config.yml` (Ignite config), `urise`/`uvrise` denominations, `sunrise` account prefix.

### Development Tools
*   **Ignite CLI**: Core tool for development and chain management.
*   **buf**: For Protobuf lifecycle management.
*   **golangci-lint**: For code linting.
*   **gomock**: For creating mock interfaces for tests.

### Technical Debt
*   There is a known issue of using deprecated Protobuf and gRPC libraries across multiple modules. A project-wide update of the Protobuf/gRPC toolchain is recommended.

This `GEMINI.md` will be updated as I learn more about the project.
