package util

import (
	"os"
	"path/filepath"
	"slices"
	"unicode"
)

func FindFile(path, fileName string) (string, error) {
	var foundPath string
	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == fileName {
			foundPath = currentPath
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", os.ErrNotExist
	}

	return foundPath, nil
}

func OnlyLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IncludesOnly(s string, runes []rune) bool {
	for _, r := range s {
		if !slices.Contains(runes, r) {
			return false
		}
	}

	return true
}

func Delete(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	return nil
}
