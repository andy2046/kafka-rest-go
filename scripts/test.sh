#!/usr/bin/env bash

set -euo pipefail

cd kafka
go fmt
go vet
golint
GOCACHE=off go test -v -race
