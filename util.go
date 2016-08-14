package fscommit

import (
	"os"
	"path/filepath"

	"github.com/dchest/uniuri"
)

func tmpPath(base string, prefix string) (string, error) {
	if len(base) <= 0 {
		base = os.TempDir()
	}
	for {
		rand := filepath.Join(base, prefix+uniuri.New())
		if _, err := os.Stat(rand); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return "", err
		}
		return rand, nil
	}
}
