package helpers

import "testing"

func TestContains(t *testing.T) {
	type Test[T comparable] struct {
		name      string
		inElement T
		inArray   []T
		out       bool
	}
	intTests := []Test[int]{
		{"int: exists", 1, []int{1, 2, 3}, true},
		{"int: does not exist", 4, []int{1, 2, 3}, false},
	}
	stringTests := []Test[string]{
		{"string: exists", "a", []string{"a", "b", "c"}, true},
		{"string: does not exist", "d", []string{"a", "b", "c"}, false},
	}

	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			if Contains(tt.inArray, tt.inElement) != tt.out {
				t.Errorf("Expected %t, got %t", tt.out, Contains(tt.inArray, tt.inElement))
			}
		})
	}

	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			if Contains(tt.inArray, tt.inElement) != tt.out {
				t.Errorf("Expected %t, got %t", tt.out, Contains(tt.inArray, tt.inElement))
			}
		})
	}
}
