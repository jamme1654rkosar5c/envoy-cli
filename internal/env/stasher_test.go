package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeStashEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestStash_SavesEntries(t *testing.T) {
	store := map[string]StashEntry{}
	entries := makeStashEntries("DB_HOST", "localhost", "DB_PORT", "5432")
	if err := Stash("backup", entries, store, DefaultStashOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := store["backup"]; !ok {
		t.Fatal("expected stash 'backup' to exist")
	}
	if len(store["backup"].Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(store["backup"].Entries))
	}
}

func TestStash_NoOverwrite_ReturnsError(t *testing.T) {
	store := map[string]StashEntry{}
	entries := makeStashEntries("KEY", "val")
	opts := DefaultStashOptions()
	_ = Stash("s1", entries, store, opts)
	if err := Stash("s1", entries, store, opts); err == nil {
		t.Fatal("expected error when overwriting without AllowOverwrite")
	}
}

func TestStash_AllowOverwrite_Succeeds(t *testing.T) {
	store := map[string]StashEntry{}
	entries := makeStashEntries("KEY", "val")
	opts := StashOptions{AllowOverwrite: true, RestoreOnPop: true}
	_ = Stash("s1", entries, store, opts)
	new := makeStashEntries("KEY", "new_val")
	if err := Stash("s1", new, store, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store["s1"].Entries[0].Value != "new_val" {
		t.Errorf("expected new_val, got %s", store["s1"].Entries[0].Value)
	}
}

func TestPop_RestoresEntries(t *testing.T) {
	store := map[string]StashEntry{}
	stashed := makeStashEntries("NEW_KEY", "new_val")
	dst := makeStashEntries("EXISTING", "exists")
	_ = Stash("s1", stashed, store, DefaultStashOptions())

	result, err := Pop("s1", dst, store, DefaultStashOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries after pop, got %d", len(result))
	}
	if _, ok := store["s1"]; ok {
		t.Error("expected stash to be removed after pop")
	}
}

func TestPop_NotFound_ReturnsError(t *testing.T) {
	store := map[string]StashEntry{}
	_, err := Pop("missing", nil, store, DefaultStashOptions())
	if err == nil {
		t.Fatal("expected error for missing stash")
	}
}

func TestPop_DoesNotOverwriteExistingKey(t *testing.T) {
	store := map[string]StashEntry{}
	stashed := makeStashEntries("DB_HOST", "remote")
	dst := makeStashEntries("DB_HOST", "localhost")
	_ = Stash("s1", stashed, store, DefaultStashOptions())

	result, err := Pop("s1", dst, store, DefaultStashOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Value != "localhost" {
		t.Errorf("expected dst value to be preserved, got %s", result[0].Value)
	}
}

func TestListStashes_ReturnsNames(t *testing.T) {
	store := map[string]StashEntry{}
	entries := makeStashEntries("K", "v")
	_ = Stash("a", entries, store, DefaultStashOptions())
	_ = Stash("b", entries, store, DefaultStashOptions())
	names := ListStashes(store)
	if len(names) != 2 {
		t.Errorf("expected 2 stash names, got %d", len(names))
	}
}
