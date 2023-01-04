package core

import "testing"

func TestCheckGlobsMatch(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		patterns []string
		out      bool
	}{
		{"empty string", "", []string{"*"}, false},
		{"empty string, no match", "", []string{"/nomatch", "./*", "~/*"}, false},
		{"empty patterns", "test", []string{}, false},
		{"no match", "/path/nomatch/file", []string{"/path/match/**/file"}, false},
		{"match", "/path/match/file", []string{"/path/match/**/file"}, true},
		{"multiple matches", "/path/match/file", []string{"/path/match/**/file", "/path/match/file"}, true},
		{"multiple matches, no match", "/path/nomatch/file", []string{"/path/match/**/file", "/path/match/file"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matches := CheckIfGlobsMatch(test.in, test.patterns...)

			if matches != test.out {
				t.Errorf("expected %v, got %v", test.out, matches)
			}
		})
	}
}

func TestFindFiles(t *testing.T) {
	t.Skip("TODO: I don't know how to test this function without creating a bunch of files and directories in the test directory. I don't want to do that.")
}
