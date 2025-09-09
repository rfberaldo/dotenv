package dotenv

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
)

var (
	override          atomic.Bool
	requireFileExists atomic.Bool
)

func SetOverride() {
	override.Store(true)
}

func SetRequireFileExists() {
	requireFileExists.Store(true)
}

// Load loads env vars from multiple sources and sets them using [os.Setenv].
//
// Load will not overwrite vars by default, first to set wins.
// To make it override use [SetOverride], when overriding, last to set wins.
//
// Load will ignore if file do not exists by default, use [SetRequireFileExists] to return error for that.
//
// Paths should be forward-slash separated, if no path is given ".env" will be used.
func Load(paths ...string) error {
	if len(paths) == 0 {
		paths = []string{".env"}
	}
	for _, path := range paths {
		if err := loadFile(path); err != nil {
			if !requireFileExists.Load() && errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}
	}
	return nil
}

// Parse reads every line from s, empty lines and lines that starts with '#' are discarded,
// the line is then split by the first '=' to get key and value.
//
// Value should not contain line breaks or quotes, all characters after the
// first '=' up to the line break are considered.
//
// Returns a key-value pair map.
func Parse(s string) (map[string]string, error) {
	return parse(strings.NewReader(s))
}

func loadFile(path string) error {
	kv, err := parseFile(path)
	if err != nil {
		return err
	}

	for k, v := range kv {
		if !override.Load() && exists(k) {
			continue
		}
		if err := os.Setenv(k, os.ExpandEnv(v)); err != nil {
			return fmt.Errorf("dotenv: setting %s=%q: %w", k, v, err)
		}
	}

	return nil
}

func exists(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}
