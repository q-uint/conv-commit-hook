// SPDX-License-Identifier: MPL-2.0

// Copyright (c) 2026 Quint Daenen.
// This file is part of conv-commit-hook.
//
// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, You can
// obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/0x51-dev/upeg/abnf/core"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
)

var logger = log.New(os.Stdout, "[CONV-COMMIT INFO] ", 0)
var debugLogger = log.New(os.Stdout, "[CONV-COMMIT DEBUG] ", 0)
var debugMode = os.Getenv("CONV_COMMIT_HOOK_DEBUG") == "true"
var errLogger = log.New(os.Stderr, "[CONV-COMMIT ERROR] ", 0)

func main() {
	if _, err := exec.LookPath("git"); err != nil {
		logger.Println("`git` not found")
		return // Skip this hook
	}
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		errLogger.Fatalf(`could not get branch: %s`, err)
	}
	branch := strings.TrimSpace(string(out))
	p, err := parser.New([]rune(branch))
	if err != nil {
		errLogger.Fatal(err)
	}

	alphaNum := append(core.ALPHA, core.DIGIT)
	special := op.Or{'-', '.'}
	noun := op.And{
		alphaNum,
		op.ZeroOrMore{Value: append(
			alphaNum,
			op.And{special, op.Peek{Value: op.Not{Value: special}}},
		)},
	}
	if _, err := p.MatchEOF(op.And{
		noun,
		op.Optional{
			Value: op.And{
				'/', noun,
			},
		},
	}); err != nil {
		if debugMode {
			debugLogger.Print(err)
		}
		errLogger.Fatalf("invalid branch name: %s", branch)
	}
}
