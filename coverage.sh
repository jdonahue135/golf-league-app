#!/bin/bash

# Run tests with coverage profiling
go test -coverprofile=coverage.out ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# remove coverage report
rm coverage.out

# Open the HTML coverage report in default web browser
open -a "Google Chrome" coverage.html