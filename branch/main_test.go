// SPDX-License-Identifier: MPL-2.0

// Copyright (c) 2026 Quint Daenen.
// This file is part of conv-commit-hook.
//
// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, You can
// obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"testing"

	"github.com/0x51-dev/upeg/abnf/core"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
)

func validateBranch(branch string) error {
	p, err := parser.New([]rune(branch))
	if err != nil {
		return err
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
	_, err = p.MatchEOF(op.And{
		noun,
		op.ZeroOrMore{
			Value: op.And{
				'/', noun,
			},
		},
	})
	return err
}

func TestBranchName(t *testing.T) {
	tests := []struct {
		name    string
		branch  string
		wantErr bool
	}{
		{"simple", "main", false},
		{"with hyphen", "feat-something", false},
		{"one slash", "feat/something", false},
		{"two slashes", "quint/feat/some-feat", false},
		{"three slashes", "org/quint/feat/some-feat", false},
		{"consecutive hyphens", "feat--bug", true},
		{"consecutive dots", "feat..bug", true},
		{"starts with hyphen", "-invalid", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBranch(tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("branch=%q: got err=%v, wantErr=%v", tt.branch, err, tt.wantErr)
			}
		})
	}
}
