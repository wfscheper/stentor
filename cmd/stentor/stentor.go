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
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/wfscheper/stentor/config"
	"github.com/wfscheper/stentor/fragment"
	"github.com/wfscheper/stentor/internal/templates"
	"github.com/wfscheper/stentor/newsfile"
	"github.com/wfscheper/stentor/release"
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
	configFile  *string
	date        time.Time
	release     *bool
	showVersion *bool
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

func (e Exec) Run() int { // nolint:gocognit // 31 > 30, but hard to see how to simplify
	// parse flags
	e, fs, err := e.parseFlags()
	if err != nil {
		if err == flag.ErrHelp {
			return succesfulExitCode
		}
		return genericExitCode
	}

	// show version and exit
	if *e.showVersion {
		e.displayVersion()
		return succesfulExitCode
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

	// determine path to config file
	if !filepath.IsAbs(*e.configFile) {
		*e.configFile = filepath.Join(e.WorkDir, *e.configFile)
	}

	// parse config file
	cfg, err := e.readConfig(*e.configFile)
	if err != nil {
		e.err.Println(err)
		return genericExitCode
	}

	fragmentFiles, err := cfg.FragmentFiles()
	if err != nil {
		e.err.Println(err)
		return genericExitCode
	}

	// parse into fragments
	var fragments []fragment.Fragment
	for _, fn := range fragmentFiles {
		f, err := fragment.New(fn)
		if err != nil {
			// log error and continue
			e.err.Printf("ignoring invalid fragment file %s: %v", fn, err)
			continue
		}
		fragments = append(fragments, f)
	}

	r, err := release.New(cfg.Repository, cfg.Markup, version, previousVersion)
	if err != nil {
		e.err.Println(err)
		return genericExitCode
	}

	// override default date
	r.Date = e.date

	r.SetSections(cfg.Sections, fragments)

	buf := &bytes.Buffer{}
	if err := generateRelease(buf, cfg, r); err != nil {
		e.err.Println(err)
		return genericExitCode
	}

	if !*e.release {
		e.out.Print(buf.String())
		return succesfulExitCode
	}

	if err := newsfile.WriteFragments(
		cfg.NewsFile,
		cfg.StartComment(),
		append([]byte("\n"), buf.Bytes()...),
		cfg.HeaderTemplate == "",
	); err != nil {
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

	e.configFile = flags.String(
		"config",
		getEnvString(e.Env, "config", filepath.Join(".stentor.d", "stentor.toml")),
		"path to config file",
	)

	date := flags.String(
		"date",
		getEnvString(e.Env, "date", time.Now().Format("2006-01-02")),
		"date of release",
	)

	e.release = flags.Bool(
		"release",
		getEnvBool(e.Env, "release", false),
		"update newsfile with fragments",
	)

	e.showVersion = flags.Bool("version", false, "show version information")

	// setup usage information
	e.setUsage(flags)

	// parse command line arguments
	err := flags.Parse(e.Args[1:])
	if err != nil {
		return e, nil, err
	}

	e.date, err = time.Parse("2006-01-02", *date)
	if err != nil {
		e.err.Println(err)
		return e, nil, err
	}

	return e, flags, err
}

func (e Exec) displayVersion() {
	e.out.Printf("%s %s built from %s on %s\n", appName, version, commit, buildDate)
}

func (Exec) readConfig(fn string) (config.Config, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return config.Config{}, fmt.Errorf("could not read config files: %w", err)
	}

	cfg, err := config.ParseBytes(data)
	if err != nil {
		return cfg, fmt.Errorf("could not parse config file: %w", err)
	}

	if err := config.ValidateConfig(cfg); err != nil {
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

func generateRelease(w io.Writer, cfg config.Config, r *release.Release) error {
	var loadTemplate = func(name, fallback string) (*template.Template, error) {
		if name != "" {
			return templates.Parse(filepath.Join(cfg.FragmentDir, name))
		}
		return templates.New(fallback)
	}

	if cfg.HeaderTemplate != "" {
		headerTemplate, err := loadTemplate(cfg.HeaderTemplate, cfg.Markup+"-header")
		if err != nil {
			return fmt.Errorf("cannot parse header template: %w", err)
		}

		if err := headerTemplate.Execute(w, r); err != nil {
			return fmt.Errorf("cannot render header template: %w", err)
		}
	}

	sectionTemplate, err := loadTemplate(cfg.SectionTemplate, cfg.Hosting+"-"+cfg.Markup+"-section")
	if err != nil {
		return fmt.Errorf("cannot parse section template: %w", err)
	}

	if err := sectionTemplate.Execute(w, r); err != nil {
		return fmt.Errorf("cannot render section template: %w", err)
	}

	return nil
}

func getEnvBool(env []string, key string, def bool) bool {
	key = strings.ToUpper(appName) + "_" + strings.ToUpper(key)
	if v, ok := lookupEnv(env, key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}

	return def
}

func getEnvString(env []string, key, def string) string {
	key = strings.ToUpper(appName) + "_" + strings.ToUpper(key)
	if v, ok := lookupEnv(env, key); ok {
		return v
	}

	return def
}

func lookupEnv(env []string, key string) (v string, ok bool) {
	for _, e := range env {
		if strings.HasPrefix(e, key+"=") {
			v = strings.Split(e, "=")[1]
			ok = true
			break
		}
	}

	return
}
