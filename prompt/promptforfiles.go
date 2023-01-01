package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	cf "github.com/osamaadam/gocfgr/configfile"
	"github.com/pkg/errors"
)

func PromptWorkAround(files []*cf.ConfigFile) (m map[string]*cf.ConfigFile, arr []string) {
	m = make(map[string]*cf.ConfigFile, len(files))
	for _, file := range files {
		readableName := file.String()
		m[readableName] = file
		arr = append(arr, readableName)
	}
	return m, arr
}

func PromptForFileSelection(files []*cf.ConfigFile) (selectedFiles []*cf.ConfigFile, err error) {
	m, arr := PromptWorkAround(files)

	prompt := &survey.MultiSelect{
		Message:  "Which files would you like to track?",
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
