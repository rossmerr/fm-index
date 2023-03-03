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
			name:    "fox",
			text:    "The quick brown fox jumps over the lazy dog",
			pattern: "do",
			offset:  []int{40},
		},
		{
			name:    "fox",
			text:    "The quick brown fox jumps over the lazy dog",
			pattern: "jumps",
			offset:  []int{20},
		},
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
		{
			name:    "abaaba",
			text:    "abaaba",
			pattern: "baab",
			offset:  []int{1},
		},
		{
			name:    "abaaba",
			text:    "abaaba",
			pattern: "abaa",
			offset:  []int{0},
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

func TestFMIndex_Extract(t *testing.T) {

	tests := []struct {
		name   string
		text   string
		offset int
		length int
		want   string
	}{
		{
			name:   "fox",
			text:   "The quick brown fox jumps over the lazy dog",
			want:   "do",
			offset: 40,
			length: 2,
		},
		{
			name:   "fox",
			text:   "The quick brown fox jumps over the lazy dog",
			want:   "jumps",
			offset: 20,
			length: 5,
		},
		{
			name:   "abaaba",
			text:   "abaaba",
			offset: 0,
			length: 3,
			want:   "aba",
		},
		{
			name:   "abaaba",
			text:   "abaaba",
			offset: 1,
			length: 3,
			want:   "baa",
		},
		{
			name:   "abaaba",
			text:   "abaaba",
			offset: 1,
			length: 4,
			want:   "baab",
		},
		{
			name:   "abaaba",
			text:   "abaaba",
			offset: 2,
			length: 4,
			want:   "aaba",
		},
		{
			name:   "abaaba",
			text:   "abaaba",
			offset: 4,
			length: 2,
			want:   "ba",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewFMIndex(tt.text)
			got := s.Extract(tt.offset, tt.length)
			if got != tt.want {
				t.Errorf("FMIndex.Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}
