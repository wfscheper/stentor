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

// +build mage

package main

import (
	"context"
	"time"

	"github.com/magefile/mage/mg"

	// mage:import
	"github.com/wfscheper/magelib"
)

const (
	moduleRice = "github.com/GeertJohan/go.rice/rice"
)

var (
	// map go:clean to clean
	Aliases = map[string]interface{}{
		"clean": magelib.Go.Clean,
	}

	// Default mage target
	Default = All

	getGolangciLint = magelib.GetGolangciLint("v1.33.0")
	getGotestsum    = magelib.GetGotestsum("v0.6.0")
	getGoreleaser   = magelib.GetGoreleaser("v0.149.0")
	getRice         = magelib.GetGoTool(moduleRice, "rice", "v1.0.0")
	getStentor      = magelib.GetStentor("v0.2.2")
)

func init() {
	// Set executable name
	magelib.ExeName = "stentor"
	// Set main package
	magelib.MainPackage = "./cmd/stentor"
	// Set test timeout
	magelib.TestTimeout = 60 * time.Second

	magelib.ChangelogDeps = []interface{}{
		func(ctx context.Context) error { return getStentor(ctx) },
	}
	magelib.GenerateDeps = []interface{}{
		func(ctx context.Context) error { return getRice(ctx) },
	}
	magelib.LintDeps = []interface{}{
		func(ctx context.Context) error { return getGolangciLint(ctx) },
	}
	magelib.ReleaseDeps = []interface{}{
		func(ctx context.Context) error { return getGoreleaser(ctx) },
	}
	magelib.TestDeps = []interface{}{
		func(ctx context.Context) error { return getGotestsum(ctx) },
	}

	magelib.ProjectTools = map[string]magelib.ToolFunc{
		magelib.ModuleGolangciLint: getGolangciLint,
		magelib.ModuleGotestsum:    getGotestsum,
		magelib.ModuleGoreleaser:   getGoreleaser,
		moduleRice:                 getRice,
		magelib.ModuleStentor:      getStentor,
	}
}

// All runs format, lint, vet, build, and test targets
func All(ctx context.Context) {
	mg.SerialCtxDeps(ctx, magelib.Go.Lint, magelib.Go.Exec, magelib.Go.Test)
}
