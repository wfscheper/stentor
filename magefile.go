// Copyright © 2020 The Stentor Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const (
	// tests
	testDir = "tests"

	// tools
	golangcilintVersion = "v1.26.0"
	gotestsumVersion    = "v0.4.1"
	riceVersion         = "v1.0.0"
	toolsDir            = "tools"
)

var (
	// Default mage target
	Default = All

	exeName = "stentor"

	goexe = "go"

	// tests
	coverageDir     = filepath.Join(testDir, "coverage")
	coverageProfile = filepath.Join(coverageDir, "coverage.out")

	// tools
	toolsBinDir      = filepath.Join(toolsDir, "bin")
	golangcilintPath = filepath.Join(toolsBinDir, "golangci-lint")
	goreleaserPath   = filepath.Join(toolsBinDir, "goreleaser")
	gotestsumPath    = filepath.Join(toolsBinDir, "gotestsum")
	ricePath         = filepath.Join(toolsBinDir, "rice")

	// commands
	gobuild = sh.RunCmd(goexe, "build")
	gofmt   = sh.RunCmd(goexe, "fmt")
	govet   = sh.RunCmd(goexe, "vet")
	rm      = sh.RunCmd("rm", "-f")

	// args
	gotestArgs = []string{"--", "-timeout=15s"}
)

func init() {
	// Force use of go modules
	os.Setenv("GO111MODULES", "on")
	if runtime.GOOS == "windows" {
		exeName += ".exe"
		golangcilintPath += ".exe"
		goreleaserPath += ".exe"
		gotestsumPath += ".exe"
	}
}

// All runs format, lint, vet, build, and test targets
func All(ctx context.Context) {
	mg.SerialCtxDeps(ctx, Lint, Vet, Build, Test)
}

// Benchmark runs the benchmark suite
func Benchmark(ctx context.Context) error {
	return runTests("-run=__absolutelynothing__", "-bench")
}

// Build runs go build
func Build(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate)
	say("building " + exeName)
	version, err := sh.Output("git", "describe", "--tags", "--always", "--dirty", "--match=v*")
	if err != nil {
		return err
	}
	commit, err := sh.Output("git", "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	buildDate := time.Now().UTC()
	ldflags := "-X main.version=" + version +
		" -X main.commit=" + commit +
		" -X main.buildDate=" + buildDate.Format(time.RFC3339)
	return gobuild("-v", "-o", filepath.Join("bin", exeName), "-ldflags", ldflags, "./cmd/stentor")
}

// Clean removes generated files
func Clean(ctx context.Context) error {
	say("cleaning files")
	return rm("-r", "bin", testDir, toolsBinDir)
}

// Coverage generates coverage reports
func Coverage(ctx context.Context) error {
	mg.CtxDeps(ctx, getGotestsum, mkCoverageDir)

	mode := os.Getenv("COVERAGE_MODE")
	if mode == "" {
		mode = "atomic"
	}
	if err := runTests(
		"-cover",
		"-covermode",
		mode,
		"-coverprofile="+coverageProfile,
	); err != nil {
		return err
	}
	if err := sh.Run(
		goexe,
		"tool",
		"cover",
		"-html="+coverageProfile,
		"-o",
		filepath.Join(coverageDir, "index.html"),
	); err != nil {
		return err
	}
	return nil
}

// Generate runs go generate
func Generate(ctx context.Context) error {
	mg.CtxDeps(ctx, getRice)

	rebuild, err := target.Dir(
		filepath.Join("internal", "templates", "rice-box.go"),
		filepath.Join("internal", "templates", "templates"),
	)
	if err == nil && rebuild {
		say("running go generate")
		return sh.RunV(goexe, "generate", "-x", "./...")
	}
	return err
}

// Lint runs golangci-lint
func Lint(ctx context.Context) error {
	mg.CtxDeps(ctx)
	say("running pre-commit hooks")
	return sh.RunV("pre-commit", "run", "--all-files")
}

// Test runs the test suite
func Test(ctx context.Context) error {
	mg.CtxDeps(ctx, getGotestsum)
	say("running tests")
	return runTests()
}

// TestRace runs the test suite with race detection
func TestRace(ctx context.Context) error {
	mg.CtxDeps(ctx, getGotestsum)
	say("running race condition tests")
	return runTests("-race")
}

// TestShort runs only tests marked as short
func TestShort(ctx context.Context) error {
	mg.CtxDeps(ctx, getGotestsum)
	say("running short tests")
	return runTests("-short")
}

// Vet runs go vet
func Vet(ctx context.Context) error {
	say("running go vet")
	return govet("./...")
}

func goGet(ctx context.Context, s string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "go", "get", s)
	cmd.Dir = filepath.Join(wd, toolsDir)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOBIN="+filepath.Join(wd, toolsBinDir))
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func getGolangciLint(ctx context.Context) error {
	rebuild, err := target.Path(golangcilintPath)
	if err == nil && rebuild {
		return goGet(ctx, "github.com/golangci/golangci-lint/cmd/golangci-lint@"+golangcilintVersion)
	}
	return err
}

func getGotestsum(ctx context.Context) error {
	rebuild, err := target.Path(gotestsumPath)
	if err == nil && rebuild {
		return goGet(ctx, "gotest.tools/gotestsum@"+gotestsumVersion)
	}
	return err
}

func getRice(ctx context.Context) error {
	rebuild, err := target.Path(ricePath)
	if err == nil && rebuild {
		return goGet(ctx, "github.com/GeertJohan/go.rice/rice@"+riceVersion)
	}
	return err
}

func mkCoverageDir(ctx context.Context) error {
	_, err := os.Stat(coverageDir)
	if os.IsNotExist(err) {
		return os.MkdirAll(coverageDir, 0755)
	}
	return err
}

func runTests(testType ...string) error {
	if update, err := strconv.ParseBool(os.Getenv("UPDATE_GOLDEN")); err == nil && update {
		testType = append(testType, "./cmd/stentor", "-update")
	} else {
		testType = append(testType, "./...")
	}
	testType = append(gotestArgs, testType...)
	return sh.RunV(gotestsumPath, testType...)
}

func say(format string, args ...interface{}) (int, error) {
	format = strings.TrimSpace(format)
	return fmt.Printf("▶ "+format+"…\n", args...)
}
