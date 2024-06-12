#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=x/liquiditypool/types/expected_keepers.go -package testutil -destination x/liquiditypool/testutil/expected_keepers_mocks.go
