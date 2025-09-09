package dotenv_test

import (
	"errors"
	"os"
	"testing"

	"github.com/rfberaldo/dotenv"
)

func TestLoad(t *testing.T) {
	t.Setenv("ENV", "dev") // predefined var

	err := dotenv.Load("testdata/.env", "testdata/.env.production")
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}

	if got, want := os.Getenv("ENV"), "dev"; got != want {
		t.Log("should not override predefined var")
		t.Errorf("expected=%s, got=%s", want, got)
	}

	if got, want := os.Getenv("HOST"), "localhost:8000"; got != want {
		t.Errorf("expected=%s, got=%s", want, got)
	}

	if got, want := os.Getenv("URL"), "http://localhost:8000"; got != want {
		t.Errorf("expected=%s, got=%s", want, got)
	}

	if got, want := os.Getenv("ANYTHING"), "what ever have after the equal sign =/\\\"`'||; works"; got != want {
		t.Errorf("expected=%s, got=%s", want, got)
	}
}

func TestLoad_Override(t *testing.T) {
	t.Setenv("ENV", "dev") // predefined var

	dotenv.SetOverride()

	err := dotenv.Load("testdata/.env", "testdata/.env.production")
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}

	if got, want := os.Getenv("ENV"), "production"; got != want {
		t.Errorf("expected=%s, got=%s", want, got)
	}

	if got, want := os.Getenv("HOST"), "example.com"; got != want {
		t.Errorf("expected=%s, got=%s", want, got)
	}

	if got, want := os.Getenv("URL"), "https://example.com"; got != want {
		t.Errorf("expected=%s, got=%s", want, got)
	}
}

func TestLoad_FileNotExists(t *testing.T) {
	if err := dotenv.Load(); err != nil {
		t.Errorf("error not expected: %s", err)
	}

	dotenv.SetRequireFileExists()
	if err := dotenv.Load(); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist: %s", err)
	}
}
