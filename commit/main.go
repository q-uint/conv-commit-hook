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
