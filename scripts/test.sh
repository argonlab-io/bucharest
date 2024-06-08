#!/bin/sh

go test ./... -coverprofile=coverage.out -json > test-report.out