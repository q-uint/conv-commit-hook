// SPDX-License-Identifier: MPL-2.0

// Copyright (c) 2026 Quint Daenen.
// This file is part of conv-commit-hook.
//
// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, You can
// obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func runCommitHook(t *testing.T, msg string) error {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "commit-hook")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %s\n%s", err, out)
	}

	msgFile := filepath.Join(t.TempDir(), "COMMIT_MSG")
	if err := os.WriteFile(msgFile, []byte(msg), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin, msgFile)
	_, err := cmd.CombinedOutput()
	return err
}

func TestCommitScope(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr bool
	}{
		{"simple scope", "feat(scope): add feature\n", false},
		{"hyphenated scope", "feat(x-y): add feature\n", false},
		{"multi-hyphen scope", "feat(foo-bar-baz): add feature\n", false},
		{"no scope", "feat: add feature\n", false},
		{"breaking", "feat!: add feature\n", false},
		{"scope with breaking", "feat(scope)!: add feature\n", false},
		{"scope leading hyphen", "feat(-scope): add feature\n", true},
		{"scope consecutive hyphens", "feat(x--y): add feature\n", true},
		{"scope trailing hyphen", "feat(scope-): add feature\n", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runCommitHook(t, tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msg=%q: got err=%v, wantErr=%v", tt.msg, err, tt.wantErr)
			}
		})
	}
}
