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

// +build mage

package main

import (
	"context"

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

	getGolangciLint = magelib.GetGolangciLint("v1.31.0")
	getGotestsum    = magelib.GetGotestsum("v0.5.3")
	getGoreleaser   = magelib.GetGoreleaser("v0.143.0")
	getRice         = magelib.GetGoTool(moduleRice, "rice", "v1.0.0")
)

func init() {
	magelib.ExeName = "stentor"
	magelib.MainPackage = "./cmd/stentor"

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
	}
}

// All runs format, lint, vet, build, and test targets
func All(ctx context.Context) {
	mg.SerialCtxDeps(ctx, magelib.Go.Lint, magelib.Go.Exec, magelib.Go.Test)
}
