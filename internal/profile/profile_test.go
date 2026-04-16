package profile_test

import (
	"os"
	"testing"

	"github.com/your-org/envoy-cli/internal/profile"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "profile-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSave_AndGet_RoundTrip(t *testing.T) {
	dir := tempDir(t)
	p := profile.Profile{Name: "local", Files: []string{".env", ".env.local"}}

	if err := profile.Save(dir, p); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := profile.Get(dir, "local")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != p.Name || len(got.Files) != len(p.Files) {
		t.Errorf("got %+v, want %+v", got, p)
	}
}

func TestSave_ReplacesExistingProfile(t *testing.T) {
	dir := tempDir(t)
	original := profile.Profile{Name: "staging", Files: []string{".env.staging"}}
	updated := profile.Profile{Name: "staging", Files: []string{".env.staging", ".env.staging.local"}}

	_ = profile.Save(dir, original)
	_ = profile.Save(dir, updated)

	got, err := profile.Get(dir, "staging")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(got.Files) != 2 {
		t.Errorf("expected 2 files after update, got %d", len(got.Files))
	}
}

func TestGet_NotFound_ReturnsError(t *testing.T) {
	dir := tempDir(t)
	_, err := profile.Get(dir, "nonexistent")
	if err == nil {
		t.Error("expected error for missing profile, got nil")
	}
}

func TestList_ReturnsAllProfiles(t *testing.T) {
	dir := tempDir(t)
	_ = profile.Save(dir, profile.Profile{Name: "dev", Files: []string{".env"}})
	_ = profile.Save(dir, profile.Profile{Name: "prod", Files: []string{".env.prod"}})

	profiles, err := profile.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(profiles))
	}
}

func TestList_EmptyDir_ReturnsEmpty(t *testing.T) {
	dir := tempDir(t)
	profiles, err := profile.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("expected 0 profiles, got %d", len(profiles))
	}
}

func TestDelete_RemovesProfile(t *testing.T) {
	dir := tempDir(t)
	_ = profile.Save(dir, profile.Profile{Name: "tmp", Files: []string{".env.tmp"}})

	if err := profile.Delete(dir, "tmp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := profile.Get(dir, "tmp")
	if err == nil {
		t.Error("expected error after deletion, got nil")
	}
}

func TestDelete_NotFound_ReturnsError(t *testing.T) {
	dir := tempDir(t)
	if err := profile.Delete(dir, "ghost"); err == nil {
		t.Error("expected error deleting non-existent profile, got nil")
	}
}

func TestSave_EmptyName_ReturnsError(t *testing.T) {
	dir := tempDir(t)
	err := profile.Save(dir, profile.Profile{Name: "", Files: []string{".env"}})
	if err == nil {
		t.Error("expected error for empty profile name, got nil")
	}
}
