package util

import "path/filepath"

type void struct{}

var null void

func UniqueFiles(files []string) (uniqueFiles []string, err error) {
	set := make(map[string]void)

	for _, file := range files {
		set[filepath.Base(file)] = null
	}

	for k := range set {
		uniqueFiles = append(uniqueFiles, k)
	}

	return
}
