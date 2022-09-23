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
	"path/filepath"
	"strings"
	"testing"

	"github.com/wfscheper/stentor/internal/test"
)

// Entry point for running integration tests.
func TestIntegration(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	var relPath = "testdata"
	err = filepath.Walk(relPath, func(path string, info os.FileInfo, err error) error {
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
			runTest(testName, relPath, wd, runMain)
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
