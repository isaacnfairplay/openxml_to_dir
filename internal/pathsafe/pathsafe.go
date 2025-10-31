package pathsafe

import (
	"errors"
	"path/filepath"
	"strings"
)

var ErrOutsideBase = errors.New("path resolves outside destination")

func Join(base, target string) (string, error) {
	cleaned := filepath.Clean(target)
	if filepath.IsAbs(cleaned) {
		return "", ErrOutsideBase
	}
	dest := filepath.Join(base, cleaned)
	ok, err := isWithin(base, dest)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrOutsideBase
	}
	return dest, nil
}

func isWithin(base, dest string) (bool, error) {
	baseAbs, err := filepath.Abs(base)
	if err != nil {
		return false, err
	}
	destAbs, err := filepath.Abs(dest)
	if err != nil {
		return false, err
	}
	rel, err := filepath.Rel(baseAbs, destAbs)
	if err != nil {
		return false, err
	}
	if rel == ".." {
		return false, nil
	}
	prefix := ".." + string(filepath.Separator)
	if strings.HasPrefix(rel, prefix) {
		return false, nil
	}
	return true, nil
}
