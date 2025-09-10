package dotenv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type pair struct {
	k string
	v string
}

// Parse reads every line from r, empty lines and lines that starts with '#' are discarded.
//
// Value should not contain line breaks or quotes, all characters after the
// first equal sign up to the line break are considered. Values are not expanded.
func Parse(r io.Reader) (map[string]string, error) {
	pairs, err := parse(r)
	if err != nil {
		return nil, err
	}
	kv := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		kv[pair.k] = pair.v
	}
	return kv, nil
}

func parseFile(path string) ([]pair, error) {
	path = filepath.FromSlash(path)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("dotenv: opening %q file: %w", path, err)
	}

	pairs, err := parse(file)
	if err != nil {
		return nil, fmt.Errorf("dotenv: parsing %q file: %w", path, err)
	}

	return pairs, file.Close()
}

func parse(r io.Reader) ([]pair, error) {
	scanner := bufio.NewScanner(r)

	pairs := make([]pair, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("dotenv: invalid line %q: missing equal sign", line)
		}
		pairs = append(pairs, pair{strings.TrimSpace(k), strings.TrimSpace(v)})
	}

	return pairs, scanner.Err()
}
