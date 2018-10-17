#!/usr/bin/env bash

set -euo pipefail

dep ensure
go build -o kubecli
