// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	toolsDir            = "tools"
)

var (
	// Default mage target
	Default = All

	exeName = "stentor"

	goexe = "go"

	// tests
	coverageDir     = filepath.Join(testDir, "coverage."+time.Now().Format("2006-01-02T15:04:05"))
	coverageProfile = filepath.Join(coverageDir, "coverage.out")

	// tools
	toolsBinDir      = filepath.Join(toolsDir, "bin")
	golangcilintPath = filepath.Join(toolsBinDir, "golangci-lint")
	goreleaserPath   = filepath.Join(toolsBinDir, "goreleaser")
	gotestsumPath    = filepath.Join(toolsBinDir, "gotestsum")

	// commands
	gobuild      = sh.RunCmd(goexe, "build")
	gofmt        = sh.RunCmd(goexe, "fmt")
	golangcilint = sh.RunCmd(golangcilintPath, "run")
	goreleaser   = sh.RunCmd(goreleaserPath)
	gotestsum    = sh.RunCmd(gotestsumPath, "--")
	govet        = sh.RunCmd(goexe, "vet")
	rm           = sh.RunCmd("rm", "-f")
)

func init() {
	// Force use of go modules
	os.Setenv("GO111MODULES", "on")
	if runtime.GOOS == "windows" {
		exeName += ".exe"
		golangcilintPath += ".exe"
		golangcilint = sh.RunCmd(golangcilintPath, "run")
		goreleaserPath += ".exe"
		goreleaser = sh.RunCmd(goreleaserPath)
		gotestsumPath += ".exe"
		gotestsum = sh.RunCmd(gotestsumPath, "--")
	}
}

// All runs format, lint, vet, build, and test targets
func All(ctx context.Context) {
	mg.SerialCtxDeps(ctx, Format, Lint, Vet, Build, Test)
}

// Benchmark runs the benchmark suite
func Benchmark(ctx context.Context) error {
	return runTests("-run=__absolutelynothing__", "-bench")
}

// Build runs go build
func Build(ctx context.Context) error {
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
		" -X main.date=" + buildDate.Format(time.RFC3339)
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

// Format runs go fmt
func Format(ctx context.Context) error {
	say("running go fmt")
	return gofmt("./...")
}

// Lint runs golangci-lint
func Lint(ctx context.Context) error {
	mg.CtxDeps(ctx, getGolangciLint)
	say("running " + golangcilintPath)
	return golangcilint()
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

func mkCoverageDir(ctx context.Context) error {
	_, err := os.Stat(coverageDir)
	if os.IsNotExist(err) {
		return os.MkdirAll(coverageDir, 0755)
	}
	return err
}

func runTests(testType ...string) error {
	testType = append(testType, "./...")
	return gotestsum(testType...)
}

func say(format string, args ...interface{}) (int, error) {
	format = strings.TrimSpace(format)
	return fmt.Printf("▶ "+format+"…\n", args...)
}
