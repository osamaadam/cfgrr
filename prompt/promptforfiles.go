package prompt

import (
	"sort"

	"github.com/AlecAivazis/survey/v2"
	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/pkg/errors"
)

// Creates a pair of a map, and a slice of strings from a slice of ConfigFiles.
func promptWorkAround(files []*cf.ConfigFile) (m map[string]*cf.ConfigFile, arr []string) {
	m = make(map[string]*cf.ConfigFile, len(files))
	for _, file := range files {
		readableName := file.String()
		m[readableName] = file
		arr = append(arr, readableName)
	}

	sort.Strings(arr)

	return m, arr
}

// Prompts the user to select files from a list of ConfigFiles.
func PromptForFileSelection(files []*cf.ConfigFile) (selectedFiles []*cf.ConfigFile, err error) {
	m, arr := promptWorkAround(files)

	prompt := &survey.MultiSelect{
		Message:  "Which files would you like to track? (this will overwrite existing files)",
		Options:  arr,
		PageSize: 10,
	}

	filteredFiles := []string{}

	if err := survey.AskOne(prompt, &filteredFiles, survey.WithKeepFilter(true)); err != nil {
		return nil, errors.WithStack(err)
	}

	for _, file := range filteredFiles {
		selectedFiles = append(selectedFiles, m[file])
	}

	return selectedFiles, nil
}
