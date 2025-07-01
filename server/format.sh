#!/bin/bash
go mod tidy

find . -name "*.go" -type f -exec gofmt -w {} \;
find . -name "*.go" -type f -exec goimports -w {} \;

echo "All Go files have been formatted."