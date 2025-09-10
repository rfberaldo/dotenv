package dotenv

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
)

var (
	override          atomic.Bool
	requireFileExists atomic.Bool
)

// SetOverride makes [Load] and [Read] override previous env vars.
func SetOverride(value bool) {
	override.Store(value)
}

// SetRequireFileExists makes [Load] and [Read] return an error if any of the given files do not exist.
func SetRequireFileExists(value bool) {
	requireFileExists.Store(value)
}

// Load loads env vars from multiple sources and sets them using [os.Setenv].
//
// Load will not overwrite vars by default, first to set wins.
// To make it override use [SetOverride], when overriding, last to set wins.
//
// Load will ignore if file do not exists by default, use [SetRequireFileExists] to force a file to exist.
//
// Paths should be forward-slash separated, if no path is given ".env" will be used.
func Load(paths ...string) error {
	return process(os.Setenv, os.Getenv, osExists, paths...)
}

// LoadTesting is like [Load], but sets vars using [t.Setenv], useful for tests.
func LoadTesting(t testing.TB, paths ...string) error {
	setenv := func(key, value string) error {
		t.Setenv(key, value)
		return nil
	}
	return process(setenv, os.Getenv, osExists, paths...)
}

// Read is like [Load], but returns a map instead of setting the environment.
// Predefined vars are only considered for expanding.
func Read(paths ...string) (map[string]string, error) {
	kv := make(map[string]string)
	setenv := func(key, value string) error {
		kv[key] = value
		return nil
	}
	getenv := func(key string) string {
		return cmp.Or(os.Getenv(key), kv[key])
	}
	exists := func(key string) bool {
		_, exists := kv[key]
		return exists
	}
	return kv, process(setenv, getenv, exists, paths...)
}

func process(
	setenv func(key, value string) error,
	getenv func(key string) string,
	exists func(key string) bool,
	paths ...string,
) error {
	if len(paths) == 0 {
		paths = []string{".env"}
	}
	for _, path := range paths {
		if err := processFile(setenv, getenv, exists, path); err != nil {
			if !requireFileExists.Load() && errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}
	}
	return nil
}

func processFile(
	setenv func(key, value string) error,
	getenv func(key string) string,
	exists func(key string) bool,
	path string,
) error {
	pairs, err := parseFile(path)
	if err != nil {
		return err
	}

	for _, pair := range pairs {
		if !override.Load() && exists(pair.k) {
			continue
		}
		if err := setenv(pair.k, os.Expand(pair.v, getenv)); err != nil {
			return fmt.Errorf("dotenv: setting %s=%q: %w", pair.k, pair.v, err)
		}
	}

	return nil
}

func osExists(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}
