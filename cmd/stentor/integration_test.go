//go:build !windows
// +build !windows

// Copyright © 2020 The Stentor Authors
//
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

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wfscheper/stentor/internal/test"
)

// The TestMain function creates a teststentor command for testing purposes and
// deletes it after the tests have been run.
// Most of this is taken from https://github.com/golang/dep and reused here.
func TestMain(m *testing.M) {
	args := []string{"build", "-o", "test" + appName + test.ExeSuffix}
	out, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "building testxkcdpwd failed: %v\n%s", err, out)
		os.Exit(2)
	}

	// Don't let these environment variables confuse the test.
	os.Unsetenv("GOPATH")
	os.Unsetenv("GIT_ALLOW_PROTOCOL")
	if home, ccacheDir := os.Getenv("HOME"), os.Getenv("CCACHE_DIR"); home != "" && ccacheDir == "" {
		// On some systems the default C compiler is ccache.
		// Setting HOME to a non-existent directory will break
		// those systems.  Set CCACHE_DIR to cope.  Issue 17668.
		os.Setenv("CCACHE_DIR", filepath.Join(home, ".ccache"))
	}
	os.Setenv("HOME", "/test-home-does-not-exist")
	if os.Getenv("GOCACHE") == "" {
		os.Setenv("GOCACHE", "off") // because $HOME is gone
	}

	r := m.Run()

	os.Remove("test" + appName + test.ExeSuffix)

	os.Exit(r)
}

// Entry point for running integration tests.
func TestIntegration(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	var relPath = "testdata"
	err = filepath.Walk(relPath, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			t.Fatal("error walking testdata:", err)
		}

		if filepath.Base(path) != "testcase.json" {
			return nil
		}

		segments := strings.Split(path, string(filepath.Separator))
		// testName is the everything after "testdata/", excluding "testcase.json"
		testName := strings.Join(segments[1:len(segments)-1], "/")
		t.Run(testName, func(t *testing.T) {
			t.Run("external", runTest(testName, relPath, wd, execCmd))
			t.Run("internal", runTest(testName, relPath, wd, runMain))
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func runTest(name, relPath, wd string, run test.RunFunc) func(t *testing.T) {
	return func(t *testing.T) {
		testCase := test.NewCase(t, filepath.Join(wd, relPath), name)

		// Skip tests
		if testCase.Skip {
			t.Skipf("skipping %s", name)
		}

		testEnv := test.NewEnvironment(t, testCase.InitialPath(), wd, run)
		defer testEnv.Cleanup()

		// force default date
		testEnv.AddEnv("STENTOR_DATE=2006-01-02")
		for _, env := range testCase.Environ {
			t.Logf("adding environment variable: %s", env)
			testEnv.AddEnv(env)
		}

		var err error
		for i, args := range testCase.Commands {
			err = testEnv.Run(appName, args)
			if err != nil && i < len(testCase.Commands)-1 {
				t.Errorf("cmd '%s' raised an unexpected error: %s", strings.Join(args, " "), err.Error())
				return
			}
		}

		if *test.UpdateGolden {
			testCase.UpdateStdout(testEnv.GetStdout())
		} else {
			testCase.CompareError(err, testEnv.GetStderr())
			testCase.CompareOutput(testEnv.GetStdout())
		}
	}
}

func execCmd(prog string, args []string, stdout, stderr io.Writer, dir string, env []string) error {
	cmd := exec.Command(prog, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Dir = dir
	cmd.Env = env
	return cmd.Run()
}

func runMain(prog string, args []string, stdout, stderr io.Writer, dir string, env []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = r
			default:
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	exc := New(dir, append([]string{prog}, args...), env, stderr, stdout)
	if exitCode := exc.Run(); exitCode != 0 {
		err = fmt.Errorf("exit status %d", exitCode)
	}

	return
}
