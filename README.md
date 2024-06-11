# Sunrise

**sunrise** is a blockchain built using Cosmos SDK and CometBFT and created with [Ignite CLI](https://ignite.com/cli).

- [cosmos/cosmos-sdk](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.50.2)
- [SunriseLayer/sunrise-core](https://github.com/SunriseLayer/sunrise-core) a fork of [cometbft/cometbft](https://github.com/cometbft/cometbft)

## Toolchain

```shell
# install go
go version

# install ignite
ignite version
# v28.3.0

apt install -y protobuf-compiler
protoc --version
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar
```

## Diagram

```ascii
                ^  +-------------------------------+  ^
                |  |                               |  |
                |  |  State-machine = Application  |  |
                |  |                               |  |   sunrise (built with Cosmos SDK)
                |  |            ^      +           |  |
                |  +----------- | ABCI | ----------+  v
Sunrise         |  |            +      v           |  ^
validator or    |  |                               |  |
full consensus  |  |           Consensus           |  |
node            |  |                               |  |
                |  +-------------------------------+  |   sunrise-core (fork of CometBFT)
                |  |                               |  |
                |  |           Networking          |  |
                |  |                               |  |
                v  +-------------------------------+  v
```

## Install

### Source

1. [Install Go](https://go.dev/doc/install) 1.22.2
1. Clone this repo
1. Install the sunrise CLI

   ```shell
   make install
   ```

### Prebuilt binary

If you'd rather not install from source, you can download a prebuilt binary from the [releases](https://github.com/SunriseLayer/sunrise/releases) page.

1. Navigate to the latest release on <https://github.com/SunriseLayer/sunrise/releases>.
1. Download the binary for your platform (e.g. `sunrise_Linux_x86_64.tar.gz`) from the **Assets** section. Tip: if you're not sure what platform you're on, you can run `uname -a` and look for the operating system (e.g. `Linux`, `Darwin`) and architecture (e.g. `x86_64`, `arm64`).
1. Extract the archive

   ```shell
   tar -xvf sunrise_Linux_x86_64.tar.gz
   ```

1. Verify the extracted binary works

   ```shell
   ./sunrised --help
   ```

1. [Optional] verify the prebuilt binary checksum. Download `checksums.txt` and then verify the checksum:

   ```shell
   sha256sum --ignore-missing --check checksums.txt
   ```

   You should see output like this:

   ```shell
   sunrise_Linux_x86_64.tar.gz: OK
   ```

See <https://docs.sunriselayer.io/node/build-node> for more information.

## Usage

```sh
# Print help
sunrised --help
```

### Environment variables

| Variable       | Explanation                        | Default value                                            | Required |
| -------------- | ---------------------------------- | -------------------------------------------------------- | -------- |
| `SUNRISE_HOME` | Home directory for the application | User home dir. [Ref](https://pkg.go.dev/os#UserHomeDir). | Optional |

### Create your own single node devnet

```sh
# Start a single node devnet
./scripts/single-node.sh

# Publish blob data to the local devnet
sunrised tx blob pay-for-blob 0x00010203040506070809 0x48656c6c6f2c20576f726c6421 \
	--chain-id private \
	--from validator \
	--keyring-backend test \
	--fees 21000uvrise \
	--yes
```

> [!NOTE]
> The sunrised binary doesn't support signing with Ledger hardware wallets on Windows and OpenBSD.

## Contributing

This repo attempts to conform to [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) so PR titles should ideally start with `fix:`, `feat:`, `build:`, `chore:`, `ci:`, `docs:`, `style:`, `refactor:`, `perf:`, or `test:` because this helps with semantic versioning and changelog generation. It is especially important to include an `!` (e.g. `feat!:`) if the PR includes a breaking change.

<!-- This repo contains multiple go modules. When using it, rename `go.work.example` to `go.work` and run `go work sync`. -->

### Tools

1. Install [golangci-lint](https://golangci-lint.run/welcome/install) 1.57.0
1. Install [markdownlint](https://github.com/DavidAnson/markdownlint) 0.39.0
1. Install [hadolint](https://github.com/hadolint/hadolint)
1. Install [yamllint](https://yamllint.readthedocs.io/en/stable/quickstart.html)
1. Install [markdown-link-check](https://github.com/tcort/markdown-link-check)
1. Install [goreleaser](https://goreleaser.com/install/)

### Helpful Commands

```sh
# Get more info on make commands.
make help

# Build the sunrised binary into the ./build directory.
make build

# Build and install the sunrised binary into the $GOPATH/bin directory.
make install

# Run tests
make test

# Format code with linters (this assumes golangci-lint and markdownlint are installed)
make fmt

# Regenerate Protobuf files (this assumes Docker is running)
make proto-gen
```
