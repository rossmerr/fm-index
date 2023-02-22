package fmindex

import (
	"reflect"
	"testing"
)

func TestFMIndex_Search(t *testing.T) {

	tests := []struct {
		name    string
		text    string
		pattern string
		count   int
	}{
		{
			name:    "abaaba",
			text:    "abaaba",
			pattern: "aba",
			count:   2,
		},
		{
			name:    "abaaba",
			text:    "abaaba",
			pattern: "baa",
			count:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewFMIndex(tt.text)
			got := s.Count(tt.pattern)
			if !reflect.DeepEqual(got, tt.count) {
				t.Errorf("FMIndex.Count() = %v, want %v", got, tt.count)
			}
		})
	}
}
