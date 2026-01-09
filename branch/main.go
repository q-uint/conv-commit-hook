// SPDX-License-Identifier: LGPL-3.0-or-later

// Copyright (c) 2026 Quint Daenen.
// This file is part of conv-commit-hook.
//
// conv-commit-hook is free software: you can redistribute it and/or modify it
// under the terms of the GNU Lesser Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// conv-commit-hook is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser Public License for
// more details.
//
// You should have received a copy of the GNU Lesser Public License along with
// conv-commit-hook. If not, see <https://www.gnu.org/licenses/>.

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
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
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
