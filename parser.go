package dotenv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// parseFile opens the file and [parse] it, path should be forward-slash separated.
func parseFile(path string) (map[string]string, error) {
	path = filepath.FromSlash(path)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("dotenv: opening %q file: %w", path, err)
	}

	kv, err := parse(file)
	if err != nil {
		return nil, fmt.Errorf("dotenv: parsing %q file: %w", path, err)
	}

	return kv, file.Close()
}

// parse reads every line from r, empty lines and lines that starts with '#' are discarded,
// the line is then split by the first '=' to get key and value.
//
// Value should not contain line breaks or quotes, all characters after the
// first '=' up to the line break are considered.
func parse(r io.Reader) (map[string]string, error) {
	scanner := bufio.NewScanner(r)

	kv := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if k, v, ok := strings.Cut(line, "="); ok {
			kv[k] = v
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return kv, nil
}
