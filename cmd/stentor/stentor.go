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

type Exec struct {
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

func New(wd string, args, env []string, err, out io.Writer) Exec {
	return Exec{
		Args:    args,
		Env:     env,
		WorkDir: wd,
		err:     log.New(err, "", 0),
		out:     log.New(out, "", 0),
	}
}

func (e Exec) Run() int {
	// parse flags
	e, _, err := e.parseFlags()
	if err != nil {
		if err == flag.ErrHelp {
			return succesfulExitCode
		}
		return genericExitCode
	}

	// show version and exit
	if e.showVersion {
		e.displayVersion()
		return succesfulExitCode
	}

	// determine path to config file
	if !filepath.IsAbs(e.configFile) {
		e.configFile = filepath.Join(e.WorkDir, e.configFile)
		e.configFile = filepath.Clean(e.configFile)
	}

	// parse config file
	_, err = e.readConfig(e.configFile)
	if err != nil {
		e.err.Println(err)
		return genericExitCode
	}

	return succesfulExitCode
}

func (e Exec) parseFlags() (Exec, *flag.FlagSet, error) {
	flags := flag.NewFlagSet(appName, flag.ContinueOnError)
	flags.SetOutput(e.err.Writer())

	flags.StringVar(&e.configFile, "config", filepath.Join(".stentor.d", "stentor.toml"), "path to config file")
	flags.BoolVar(&e.showVersion, "version", false, "show version information")

	// setup usage information
	e.setUsage(flags)

	// parse command line arguments
	if err := flags.Parse(e.Args[1:]); err != nil {
		return e, nil, err
	}

	return e, flags, nil
}

func (e Exec) displayVersion() {
	e.out.Printf("%s %s built from %s on %s\n", appName, version, commit, buildDate)
}

func (Exec) readConfig(fn string) (Config, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return Config{}, fmt.Errorf("could not read config files: %w", err)
	}

	cfg, err := ParseBytes(data)
	if err != nil {
		return cfg, fmt.Errorf("could not parse config file: %w", err)
	}

	if err := ValidateConfig(cfg); err != nil {
		return cfg, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (e Exec) setUsage(fs *flag.FlagSet) {
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
		e.out.Printf(`Usage: %[1]s [OPTIONS]

%[1]s is a CLI for generating a change log or release notes from a set of fragment files and templates.

Flags:

%s`, appName, flagsUsage.String())
	}
}
