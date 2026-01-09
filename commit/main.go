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

	"github.com/0x51-dev/upeg/abnf/core"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
)

var debugLogger = log.New(os.Stdout, "[CONV-COMMIT DEBUG] ", 0)
var debugMode = os.Getenv("CONV_COMMIT_HOOK_DEBUG") == "true"
var errLogger = log.New(os.Stderr, "[CONV-COMMIT ERROR] ", 0)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		errLogger.Fatal("Wrong number of args passed to hook, is it installed correctly?")
	}
	commitMessage, err := os.ReadFile(args[0])
	if err != nil {
		errLogger.Fatal(err)
	}
	p, err := parser.New([]rune(string(commitMessage)))
	if err != nil {
		errLogger.Fatal(err)
	}

	// (1)
	noun := op.OneOrMore{Value: core.ALPHA}
	if _, err := p.Match(noun); err != nil {
		if debugMode {
			debugLogger.Print(err)
		}
		errLogger.Fatal("(1) Commit does NOT start with a noun.")
	}

	// (4)
	if _, err := p.Match(op.And{
		op.Optional{Value: op.And{'(', noun, ')'}},
		op.Optional{Value: '!'},
		op.Peek{Value: ":"}, // Check up until the terminal colon.
	}); err != nil {
		if debugMode {
			debugLogger.Print(err)
		}
		errLogger.Fatal("(4) Invalid scope, should be `(ALPHA*)`.")
	}

	// (1)
	if _, err := p.Match(": "); err != nil {
		if debugMode {
			debugLogger.Print(err)
		}
		errLogger.Fatal("(1) Commit does NOT end with a terminal colon and space.")
	}

	// (5)
	if _, err := p.Match(op.And{
		op.Not{Value: op.And{
			op.ZeroOrMore{Value: core.WSP},
			op.Or{core.CRLF, op.EOF{}}, // Only whitespace before the newline.
		}},
		op.OneOrMore{Value: op.AnyBut{Value: core.CRLF}},
		op.Or{core.CRLF, op.EOF{}},
	}); err != nil {
		if debugMode {
			debugLogger.Print(err)
		}
		errLogger.Fatal("(5) A description MUST immediately follow the colon and space.")
	}
}
