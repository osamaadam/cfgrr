package helpers

import (
	"bufio"
	"os"
)

// ReadFileLines reads a file line by line and returns a slice of strings.
func ReadFileLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return
}
