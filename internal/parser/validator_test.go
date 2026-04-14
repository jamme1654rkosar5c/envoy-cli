package parser

import (
	"testing"
)

func makeEnvFile(entries []EnvEntry) *EnvFile {
	return &EnvFile{Path: "test.env", Entries: entries}
}

func TestValidate_ValidEntries(t *testing.T) {
	env := makeEnvFile([]EnvEntry{
		{Key: "APP_ENV", Value: "production", Line: 1},
		{Key: "PORT", Value: "8080", Line: 2},
	})

	errs := Validate(env)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func TestValidate_InvalidKeyFormat(t *testing.T) {
	env := makeEnvFile([]EnvEntry{
		{Key: "lower_case", Value: "value", Line: 1},
	})

	errs := Validate(env)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Key != "lower_case" {
		t.Errorf("unexpected error key: %s", errs[0].Key)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	env := makeEnvFile([]EnvEntry{
		{Key: "SECRET", Value: "", Line: 1},
	})

	errs := Validate(env)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Message != "value is empty" {
		t.Errorf("unexpected error message: %s", errs[0].Message)
	}
}

func TestValidate_DuplicateKey(t *testing.T) {
	env := makeEnvFile([]EnvEntry{
		{Key: "DB_HOST", Value: "localhost", Line: 1},
		{Key: "DB_HOST", Value: "remotehost", Line: 5},
	})

	errs := Validate(env)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if errs[0].Line != 5 {
		t.Errorf("expected error on line 5, got line %d", errs[0].Line)
	}
}

func TestValidationError_String(t *testing.T) {
	e := ValidationError{Line: 3, Key: "FOO", Message: "value is empty"}
	expected := "line 3 [FOO]: value is empty"
	if e.Error() != expected {
		t.Errorf("expected %q, got %q", expected, e.Error())
	}
}
