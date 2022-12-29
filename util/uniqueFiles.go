package util

import (
	"path/filepath"
)

func UniqueFiles(files []string) (uniqueFiles []string, err error) {
	set := make(map[string]uint8)

	for _, file := range files {
		fileBase := filepath.Base(file)
		set[fileBase] += 1
	}

	i := 0

	for k, v := range set {
		if v < 2 {
			uniqueFiles = append(uniqueFiles, k)
		}
		i++
	}

	return uniqueFiles, nil
}
