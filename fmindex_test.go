package fmindex

import (
	"reflect"
	"testing"
)

func TestFMIndex_Count(t *testing.T) {

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

func TestFMIndex_Locate(t *testing.T) {

	tests := []struct {
		name    string
		text    string
		pattern string
		offset  []int
	}{
		{
			name:    "abaaba",
			text:    "abaaba",
			pattern: "aba",
			offset:  []int{3, 0},
		},
		{
			name:    "abaaba",
			text:    "abaaba",
			pattern: "baa",
			offset:  []int{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewFMIndex(tt.text)
			got := s.Locate(tt.pattern)
			if !reflect.DeepEqual(got, tt.offset) {
				t.Errorf("FMIndex.Locate() = %v, want %v", got, tt.offset)
			}
		})
	}
}
