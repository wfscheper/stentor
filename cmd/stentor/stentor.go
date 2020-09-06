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
	"io/ioutil"
	"log"
	"path/filepath"
	"text/tabwriter"
)

const (
	appName           = "stentor"
	genericExitCode   = 1
	succesfulExitCode = 0
)

var (
	buildDate = "unknown"
	commit    = "unknown"
	version   = "dev"
)

type Stentor struct {
	Args    []string // command-line arguments
	Env     []string // os environment
	WorkDir string   // Where to execute

	// internal loggers
	err *log.Logger
	out *log.Logger

	// command-line options
	configFile  string
	showVersion bool
}

func New(wd string, args, env []string, err, out io.Writer) *Stentor {
	return &Stentor{
		Args:    args,
		Env:     env,
		WorkDir: wd,
		err:     log.New(err, "", 0),
		out:     log.New(out, "", 0),
	}
}

func (s *Stentor) Run() int {
	// parse flags
	if err := s.parseFlags(); err != nil {
		if err == flag.ErrHelp {
			return succesfulExitCode
		}
		return genericExitCode
	}

	// show version and exit
	if s.showVersion {
		s.displayVersion()
		return succesfulExitCode
	}

	// determine path to config file
	if !filepath.IsAbs(s.configFile) {
		s.configFile = filepath.Join(s.WorkDir, s.configFile)
		s.configFile = filepath.Clean(s.configFile)
	}

	// parse config file
	data, err := ioutil.ReadFile(s.configFile)
	if err != nil {
		s.err.Printf("could not read config files: %v", err)
		return genericExitCode
	}

	cfg, err := ParseBytes(data)
	if err != nil {
		s.err.Printf("could not parse config file: %v", err)
		return genericExitCode
	}

	if err := ValidateConfig(cfg); err != nil {
		s.err.Printf("invalid configuration: %v", err)
		return genericExitCode
	}

	return succesfulExitCode
}

func (s *Stentor) parseFlags() error {
	flags := flag.NewFlagSet(appName, flag.ContinueOnError)
	flags.SetOutput(s.err.Writer())

	flags.StringVar(&s.configFile, "config", filepath.Join(".stentor.d", "stentor.toml"), "path to config file")
	flags.BoolVar(&s.showVersion, "version", false, "show version information")

	// setup usage information
	s.setUsage(flags)

	// parse command line arguments
	if err := flags.Parse(s.Args[1:]); err != nil {
		return err
	}

	return nil
}

func (s *Stentor) displayVersion() {
	s.out.Printf("%s %s built from %s on %s\n", appName, version, commit, buildDate)
}

func (s *Stentor) setUsage(fs *flag.FlagSet) {
	var flagsUsage bytes.Buffer
	tw := tabwriter.NewWriter(&flagsUsage, 0, 4, 2, ' ', 0)
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
		s.out.Printf(`Usage: %[1]s [OPTIONS]

%[1]s is a CLI for generating a change log or release notes from a set of fragment files and templates.

Flags:

%s`, appName, flagsUsage.String())
	}
}
