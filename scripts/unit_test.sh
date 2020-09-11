#!/bin/bash
go test -race $(go list ./... | grep -v /vendor/) -v -coverprofile=coverage.out || exit 1
go tool cover -func=coverage.out
