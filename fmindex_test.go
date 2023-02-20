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
		want    *FMIndexResult
	}{
		{
			name:    "abaaba",
			text:    "abaaba",
			pattern: "aba",
			want: &FMIndexResult{
				count: 2,
				rows: []*FMIndexResultRow{
					{
						start: 1,
						end:   3,
					},
					{
						start: 3,
						end:   3,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewFMIndex(tt.text, false)
			got := s.Search(tt.pattern)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FMIndex.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
