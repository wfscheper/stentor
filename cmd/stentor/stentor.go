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

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"runtime"
	"text/tabwriter"
)

const (
	// output strings
	versionInfo = `%s
  version     : %s
  build date  : %s
  git hash    : %s
  go version  : %s
  go compiler : %s
  platform    : %s/%s
`
)

type Stentor struct {
	Args    []string
	Env     []string
	Stderr  *log.Logger
	Stdout  *log.Logger
	WorkDir string

	showVersion *bool
}

func New(wd string, args, env []string, stderr, stdout io.Writer) *Stentor {
	return &Stentor{
		Args:    args,
		Env:     env,
		Stderr:  log.New(stderr, "", 0),
		Stdout:  log.New(stdout, "", 0),
		WorkDir: wd,
	}
}

func (s *Stentor) Run() int {
	// parse flags
	if err := s.parseFlags(); err != nil {
		return genericErrorCode
	}

	// display version info
	if *s.showVersion {
		return s.displayVersion()
	}
	return 0
}

func (s *Stentor) displayVersion() int {
	s.Stdout.Printf(
		versionInfo,
		appName,
		version,
		date,
		commit,
		runtime.Version(),
		runtime.Compiler,
		runtime.GOOS,
		runtime.GOARCH,
	)
	return successExitCode
}

func (s *Stentor) parseFlags() error {
	flags := flag.NewFlagSet(appName, flag.ContinueOnError)
	flags.SetOutput(s.Stderr.Writer())

	s.showVersion = flags.Bool("version", false, "show version information")

	s.setUsage(flags)

	if err := flags.Parse(s.Args[1:]); err != nil && err != flag.ErrHelp {
		return err
	}

	return nil
}

func (s *Stentor) setUsage(fs *flag.FlagSet) {
	var usage bytes.Buffer
	tw := tabwriter.NewWriter(&usage, 0, 4, 2, ' ', 0)
	fs.VisitAll(func(f *flag.Flag) {
		switch f.DefValue {
		case "":
			fmt.Fprintf(tw, "\t-%s\t%s\n", f.Name, f.Usage)
		case " ":
			fmt.Fprintf(tw, "\t-%s\t%s (default: '%s')\n", f.Name, f.Usage, f.DefValue)
		default:
			fmt.Fprintf(tw, "\t-%s\t%s (default: %s)\n", f.Name, f.Usage, f.DefValue)
		}
	})

	tw.Flush()

	fs.Usage = func() {
		s.Stdout.Printf(`Usage: %s [OPTIONS]

%[1]s is a CLI for generating a change log or release notes from a set of fragment files and templates.

Flags:

%s`, appName, usage.String())
	}
}
