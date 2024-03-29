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

package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Case loads a testdata.json test configuration and executes that test.
type Case struct {
	t           *testing.T
	name        string
	rootPath    string
	initialPath string
	Commands    [][]string `json:"commands"`
	Skip        bool       `json:"skip"`
	Environ     []string   `json:"environ"`
}

// NewCase returns a Case.
func NewCase(t *testing.T, dir, name string) *Case {
	rootPath := filepath.FromSlash(filepath.Join(dir, name))

	c := &Case{
		t:           t,
		name:        name,
		rootPath:    rootPath,
		initialPath: filepath.Join(rootPath, "initial"),
	}

	data, err := os.ReadFile(filepath.Join(rootPath, "testcase.json"))
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		t.Fatal(err)
	}

	return c
}

// CompareOutput compares stdout to the contents of a stdout file in the test directory.
func (c *Case) CompareOutput(stdout string) {
	data, err := os.ReadFile(filepath.Join(c.rootPath, "stdout"))
	if err != nil {
		if os.IsNotExist(err) && stdout == "" {
			return
		}
		panic(err)
	}

	assert.Equal(c.t, string(data), stdout)
}

// CompareError compares stderr to the contents of a stderr file in the test directory.
func (c *Case) CompareError(errIn error, stderr string) {
	data, err := os.ReadFile(filepath.Join(c.rootPath, "stderr"))
	if err != nil {
		if os.IsNotExist(err) && errIn == nil && stderr == "" {
			return
		}
		switch {
		case stderr != "":
			panic(errors.New(stderr))
		case errIn != nil:
			panic(errIn)
		default:
			panic(err)
		}
	}

	want := string(data)
	expectError := want != ""
	gotError := stderr != "" && errIn != nil
	switch {
	case expectError && gotError:
		assert.Regexp(c.t, want, stderr, "stderr did not match the expected error")
	case expectError && !gotError:
		c.t.Errorf("expected error:\n%s", want)
	case !expectError && gotError:
		c.t.Errorf("unexpected error:\n%s", stderr)
	}
}

func (c *Case) InitialPath() string {
	return c.initialPath
}

// UpdateStdout updates the golden file for stdout with the working result.
func (c *Case) UpdateStdout(stdout string) {
	stdoutPath := filepath.Join(c.rootPath, "stdout")
	_, err := os.Stat(stdoutPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Don't update the stdout file if it doesn't exist.
			return
		}
		panic(err)
	}

	if err := os.WriteFile(stdoutPath, []byte(stdout), 0644); err != nil {
		c.t.Fatal(err)
	}
}

// Environment defines a test execution environment and captures the output.
type Environment struct {
	t      *testing.T
	tmpdir string
	wd     string
	env    []string
	stdout bytes.Buffer
	stderr bytes.Buffer
	run    RunFunc
}

// NewEnvironment initializes the test Environment.
func NewEnvironment(t *testing.T, rootPath, wd string, run RunFunc) *Environment {
	e := &Environment{
		t:   t,
		wd:  wd,
		env: os.Environ(),
		run: run,
	}

	e.makeTempDir()
	e.CopyTree(rootPath)

	if err := os.Chdir(e.tmpdir); err != nil {
		t.Fatalf("could not cd to %s: %v", e.tmpdir, err)
	}

	return e
}

func (te *Environment) AddEnv(e string) {
	te.env = append(te.env, e)
}

func (te *Environment) Cleanup() {
	_ = os.Chdir(te.wd)
	_ = os.RemoveAll(te.tmpdir)
}

func (te *Environment) CopyTree(src string) {
	err := filepath.Walk(src, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if p != src {
			// strip src from path
			localpath := p[len(src)+1:]
			if info.IsDir() {
				te.MakeDir(localpath)
			} else {
				dst := filepath.Join(te.tmpdir, localpath)
				copyFile(te.t, dst, p)
			}
		}

		return nil
	})

	if err != nil && !os.IsNotExist(err) {
		te.t.Fatalf("could not copy %s: %v", src, err)
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

// Join returns a path rooted at the environment's tempdir.
func (te *Environment) Join(args ...string) string {
	return filepath.Join(te.tmpdir, filepath.Join(args...))
}

func (te *Environment) MakeDir(args ...string) {
	p := te.Join(args...)
	if err := os.MkdirAll(p, 0755); err != nil {
		te.t.Fatalf("could not create directories %s: %+v", p, err)
	}
}

// Run runs the tests command with args.
func (te *Environment) Run(progName string, args []string) error {
	if *Verbose {
		te.t.Logf("running `test%s %s`", progName, strings.Join(args, " "))
	}
	prog := filepath.Join(te.wd, "test"+progName+ExeSuffix)
	te.stdout.Reset()
	te.stderr.Reset()

	status := te.run(prog, args, &te.stdout, &te.stderr, te.tmpdir, te.env)

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

func (te *Environment) makeTempDir() {
	if te.tmpdir == "" {
		te.tmpdir = te.t.TempDir()

		// OSX uses a symlink, so resolve the link
		if runtime.GOOS == "darwin" {
			realDir, err := filepath.EvalSymlinks(te.tmpdir)
			if err != nil {
				te.t.Fatal(err)
			}
			te.tmpdir = realDir
		}
	}
}

// RunFunc is a function that runs a test.
type RunFunc func(prog string, args []string, stdout, stderr io.Writer, dir string, env []string) error

func copyFile(t *testing.T, dst, src string) {
	t.Helper()

	s, err := os.Open(src)
	if err != nil {
		t.Fatalf("could not open %s: %v", src, err)
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		t.Fatalf("could not create %s: %v", dst, err)
	}
	defer d.Close()

	if _, err := io.Copy(d, s); err != nil {
		t.Fatalf("could not write to %s: %v", dst, err)
	}
}
