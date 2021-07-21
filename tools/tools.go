// +build tools

package main

import (
	_ "github.com/GeertJohan/go.rice/rice"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/sassoftware/gotagger"
	_ "github.com/wfscheper/stentor"
	_ "gotest.tools/gotestsum"
)
