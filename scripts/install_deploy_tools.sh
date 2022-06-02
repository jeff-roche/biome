#!/bin/bash

# Install svu for version comparison
go install github.com/caarlos0/svu@v1.9.0

# Install goreleaser to do the release
go install github.com/goreleaser/goreleaser@v1.9.2