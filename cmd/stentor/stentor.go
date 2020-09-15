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
	"os"
	"path/filepath"
	"text/tabwriter"
	"text/template"

	"github.com/wfscheper/stentor/fragment"
	"github.com/wfscheper/stentor/internal/templates"
	"github.com/wfscheper/stentor/newsfile"
	"github.com/wfscheper/stentor/release"
	"github.com/wfscheper/stentor/section"
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
	release     bool
	showVersion bool
}

func New(wd string, args, env []string, err, out io.Writer) Exec {
	return Exec{
		Args:    args,
		Env:     env,
		WorkDir: wd,
		err:     log.New(err, appName+": ", 0),
		out:     log.New(out, "", 0),
	}
}

func (e Exec) Run() int {
	// parse flags
	e, fs, err := e.parseFlags()
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
	cfg, err := e.readConfig(e.configFile)
	if err != nil {
		e.err.Println(err)
		return genericExitCode
	}

	if len(fs.Args()) > 2 {
		e.err.Println("too many arguments")
		return genericExitCode
	}

	version := fs.Arg(0)
	if version == "" {
		e.err.Println("missing NEW argument")
		return genericExitCode
	}

	previousVersion := fs.Arg(1)
	if previousVersion == "" {
		e.err.Println("missing PREVIOUS argument")
		return genericExitCode
	}

	fragmentFiles, err := cfg.FragmentFiles()
	if err != nil {
		e.err.Println(err)
		return genericExitCode
	}

	// if there are no fragments, and no show always sections then exit
	if len(fragmentFiles) == 0 {
		var exit bool = true
		for _, s := range cfg.Sections {
			if s.ShowAlways != nil && *s.ShowAlways {
				exit = false
				break
			}
		}
		if exit {
			return succesfulExitCode
		}
	}

	r := release.New(cfg.Repository, cfg.Markup, version, previousVersion)
	r, err = configureSections(r, cfg.Sections, fragmentFiles)
	if err != nil {
		e.err.Println(err)
		return genericExitCode
	}

	buf := &bytes.Buffer{}
	if err := generateRelease(buf, cfg, r); err != nil {
		e.err.Println(err)
		return genericExitCode
	}

	if !e.release {
		e.out.Print(buf.String())
		return succesfulExitCode
	}

	if err := newsfile.WriteFragments(cfg.NewsFile, cfg.StartComment(), buf.Bytes()); err != nil {
		e.err.Printf("cannot update %s: %v", cfg.NewsFile, err)
		return genericExitCode
	}

	var failed bool
	for _, f := range fragmentFiles {
		if err := os.Remove(f); err != nil {
			e.err.Printf("cannot remove fragment file %s: %v", f, err)
			failed = true
		}
	}

	if failed {
		return genericExitCode
	}
	return succesfulExitCode
}

func (e Exec) parseFlags() (Exec, *flag.FlagSet, error) {
	flags := flag.NewFlagSet(appName, flag.ContinueOnError)
	flags.SetOutput(e.err.Writer())

	flags.StringVar(&e.configFile, "config", filepath.Join(".stentor.d", "stentor.toml"), "path to config file")
	flags.BoolVar(&e.release, "release", false, "update newsfile with fragments")
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
		e.out.Printf(`Usage: %[1]s [OPTIONS] NEW PREVIOUS

Update a news file with the changes from version PREVIOUS to NEW.

Flags:

%s`, appName, flagsUsage.String())
	}
}

func configureSections(r release.Release, sections []SectionConfig, fragmentFiles []string) (release.Release, error) {
	sectionMap := map[string]section.Section{}
	for _, fragmentFile := range fragmentFiles {
		f, section, err := fragment.New(fragmentFile)
		if err != nil {
			return r, err
		}

		s := sectionMap[section]
		s.Fragments = append(s.Fragments, f)
		sectionMap[section] = s
	}

	for _, cfg := range sections {
		if s, ok := sectionMap[cfg.ShortName]; ok {
			if cfg.ShowAlways != nil {
				s.ShowAlways = *cfg.ShowAlways
			}
			s.Title = cfg.Name
			r.Sections = append(r.Sections, s)
		} else if cfg.ShowAlways != nil && *cfg.ShowAlways {
			r.Sections = append(r.Sections, section.Section{
				ShowAlways: *cfg.ShowAlways,
				Title:      cfg.Name,
			})
		}
	}

	return r, nil
}

func generateRelease(w io.Writer, cfg Config, r release.Release) error {
	var loadTemplate = func(name, fallback string) (*template.Template, error) {
		if name != "" {
			return templates.Parse(name)
		}
		return templates.New(fallback)
	}

	headerTemplate, err := loadTemplate(cfg.HeaderTemplate, cfg.Hosting+"-"+cfg.Markup+"-header")
	if err != nil {
		return fmt.Errorf("cannot parse header template: %w", err)
	}

	sectionTemplate, err := loadTemplate(cfg.SectionTemplate, cfg.Hosting+"-"+cfg.Markup+"-section")
	if err != nil {
		return fmt.Errorf("cannot parse section template: %w", err)
	}

	if err := headerTemplate.Execute(w, r); err != nil {
		return fmt.Errorf("cannot render header template: %w", err)
	}

	if err := sectionTemplate.Execute(w, r); err != nil {
		return fmt.Errorf("cannot render section template: %w", err)
	}

	return nil
}
