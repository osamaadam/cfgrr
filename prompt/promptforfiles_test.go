package prompt

import (
	"testing"

	cf "github.com/osamaadam/cfgrr/configfile"
)

func TestPromptWorkaround(t *testing.T) {
	tests := []struct {
		name   string
		in     []*cf.ConfigFile
		outMap map[string]*cf.ConfigFile
		outArr []string
	}{
		{"empty", nil, map[string]*cf.ConfigFile{}, []string{}},
		{"one", []*cf.ConfigFile{{Path: "/home/file"}}, map[string]*cf.ConfigFile{"file - /home/file": {Path: "/home/file"}}, []string{"file - /home/file"}},
		{"multiple",
			[]*cf.ConfigFile{{Path: "/home/file"}, {Path: "/home/file2"}},
			map[string]*cf.ConfigFile{"file - /home/file": {Path: "/home/file"}, "file2 - /home/file2": {Path: "/home/file2"}},
			[]string{"file - /home/file", "file2 - /home/file2"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outMap, outArr := promptWorkaround(test.in)
			if len(outMap) != len(test.outMap) {
				t.Errorf("expected %v, got %v", test.outMap, outMap)
			}
			if len(outArr) != len(test.outArr) {
				t.Errorf("expected %v, got %v", test.outArr, outArr)
			}
		})
	}
}
