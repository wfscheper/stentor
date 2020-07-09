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

package test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ianbruene/go-difflib/difflib"
)

// Case loads a testdata.json test configuration and executes that test.
type Case struct {
	t        *testing.T
	name     string
	rootPath string
	Commands [][]string `json:"commands"`
	Skip     bool       `json:"skip"`
}

// NewCase returns a Case.
func NewCase(t *testing.T, dir, name string) *Case {
	rootPath := filepath.FromSlash(filepath.Join(dir, name))

	c := &Case{
		t:        t,
		name:     name,
		rootPath: rootPath,
	}

	data, err := ioutil.ReadFile(filepath.Join(rootPath, "testcase.json"))
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		t.Fatal(err)
	}

	return c
}

// CompareOutput compares stdout to the contents of a stdout.txt file in the test directory.
func (c *Case) CompareOutput(stdout string) {
	data, err := ioutil.ReadFile(filepath.Join(c.rootPath, "stdout.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}

	if got, want := stdout, string(data); got != want {
		c.t.Errorf("stdout was not as expected\n%s", diffErr(c.t, got, want))
	}
}

// CompareError compares stderr to the contents of a stderr.txt file in the test directory.
func (c *Case) CompareError(errIn error, stderr string) {
	data, err := ioutil.ReadFile(filepath.Join(c.rootPath, "stderr.txt"))
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}

	want := string(data)
	expectError := data != nil
	gotError := stderr != "" && errIn != nil
	switch {
	case expectError && gotError:
		switch matches := strings.Count(stderr, want); matches {
		case 0:
			c.t.Errorf("stderr did not contain the expected error:\n%s", diffErr(c.t, stderr, want))
		case 1:
		default:
			c.t.Errorf("expected error '%s' occurred %d times in stderr\n%s", want, matches, stderr)
		}
	case expectError && !gotError:
		c.t.Errorf("expected error:\n%s", want)
	case !expectError && gotError:
		c.t.Errorf("unexpected error:\n%s", stderr)
	}
}

// UpdateStderr updates the golden file for stderr with the working result.
func (c *Case) UpdateStderr(stderr string) {
	stderrPath := filepath.Join(c.rootPath, "stderr.txt")
	_, err := os.Stat(stderrPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Don't update the stdout.txt file if it doesn't exist.
			return
		}
		panic(err)
	}

	if err := ioutil.WriteFile(stderrPath, []byte(stderr), 0644); err != nil {
		c.t.Fatal(err)
	}
}

// UpdateStdout updates the golden file for stdout with the working result.
func (c *Case) UpdateStdout(stdout string) {
	stdoutPath := filepath.Join(c.rootPath, "stdout.txt")
	_, err := os.Stat(stdoutPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Don't update the stdout.txt file if it doesn't exist.
			return
		}
		panic(err)
	}

	if err := ioutil.WriteFile(stdoutPath, []byte(stdout), 0644); err != nil {
		c.t.Fatal(err)
	}
}

// Environment defines a test execution environment and captures the output.
type Environment struct {
	t      *testing.T
	wd     string
	env    []string
	stdout bytes.Buffer
	stderr bytes.Buffer
	run    RunFunc
}

// NewEnvironment initializes the test Environment.
func NewEnvironment(t *testing.T, wd string, run RunFunc) *Environment {
	return &Environment{
		t:   t,
		wd:  wd,
		env: os.Environ(),
		run: run,
	}
}

// GetStdout returns the captures stdout.
func (te *Environment) GetStdout() string {
	return te.stdout.String()
}

// GetStderr returns the captures stderr.
func (te *Environment) GetStderr() string {
	return te.stderr.String()
}

// Run runs the tests command with args.
func (te *Environment) Run(progName string, args []string) error {
	if *Verbose {
		te.t.Logf("running `test%s %s`", progName, strings.Join(args, " "))
	}
	prog := filepath.Join(te.wd, "test"+progName+ExeSuffix)
	te.stdout.Reset()
	te.stderr.Reset()

	status := te.run(prog, args, &te.stdout, &te.stderr, te.wd, te.env)

	if *Verbose {
		if te.stdout.Len() > 0 {
			te.t.Logf("\nstdout: %s", te.stdout.String())
		}
		if te.stderr.Len() > 0 {
			te.t.Logf("\nstderr: %s", te.stderr.String())
		}
	}
	return status
}

// RunFunc is a function that runs a test.
type RunFunc func(prog string, args []string, stdout, stderr io.Writer, dir string, env []string) error

func diffErr(t *testing.T, got, want string) string {
	t.Helper()

	diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(want),
		B:        difflib.SplitLines(got),
		Context:  3,
		FromFile: "want",
		ToFile:   "got",
	})
	if err != nil {
		t.Fatal(err)
	}

	return diff
}
