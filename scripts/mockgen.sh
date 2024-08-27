#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=x/liquiditypool/types/expected_keepers.go -package testutil -destination x/liquiditypool/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/liquidityincentive/types/expected_keepers.go -package testutil -destination x/liquidityincentive/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/swap/types/expected_keepers.go -package testutil -destination x/swap/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/da/types/expected_keepers.go -package testutil -destination x/da/testutil/expected_keepers_mocks.go
