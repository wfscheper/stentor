//go:build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/sassoftware/gotagger"
	_ "gotest.tools/gotestsum"
)
