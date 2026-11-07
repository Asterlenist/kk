package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendTrimsOldEntriesToLimit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.log")
	const limit = 512
	for i := 0; i < 10; i++ {
		if err := Append(path, limit, strings.Repeat("command ", 8)+string(rune('a'+i)), nil); err != nil {
			t.Fatal(err)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) > limit {
		t.Fatalf("history size = %d, want <= %d", len(data), limit)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 0 {
		t.Fatal("expected at least one history entry")
	}
	var last entry
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &last); err != nil {
		t.Fatalf("last entry is invalid JSON: %v", err)
	}
	if !strings.HasSuffix(last.Command, "j") {
		t.Fatalf("last command = %q, want newest command", last.Command)
	}
}

func TestAppendTruncatesSingleOversizedEntry(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.log")
	const limit = 128
	if err := Append(path, limit, strings.Repeat("long-command ", 100), nil); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) > limit {
		t.Fatalf("history size = %d, want <= %d", len(data), limit)
	}
	var got entry
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(data))), &got); err != nil {
		t.Fatalf("entry is invalid JSON: %v", err)
	}
	if !strings.HasSuffix(got.Command, "...") {
		t.Fatalf("oversized command was not marked as truncated: %q", got.Command)
	}
}
