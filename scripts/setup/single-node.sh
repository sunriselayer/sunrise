#!/bin/sh

# Stop script execution if an error is encountered
set -o errexit
# Stop script execution if an undefined variable is used
set -o nounset

SCRIPT_DIR=$(cd $(dirname $0); pwd)
. $SCRIPT_DIR/init.sh

go run ./cmd/sunrised/main.go start --home ./data/test