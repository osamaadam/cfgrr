package helpers

import (
	"bufio"
	"os"

	"github.com/pkg/errors"
)

// ReadFileLines reads a file line by line and returns a slice of strings.
func ReadFileLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return
}
