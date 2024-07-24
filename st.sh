#!/bin/bash
clc -s -e ufunc_test.go
cat Version.dat
go mod tidy
go fmt .
echo can\'t run staticcheck due to rangefuncs
#staticcheck .
go vet .
echo can\'t run golangci-lint due to rangefuncs
#golangci-lint run
git st
