#!/bin/bash

current_dir="$(pwd)"

echo "current_dir: $current_dir"

rm -rf "$current_dir/../data"

echo "data directory removed"

go run -race "$current_dir/../cmd/main.go"
