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
		offset  []int
	}{
		// {
		// 	name:    "abaaba",
		// 	text:    "abaaba",
		// 	pattern: "aba",
		// 	count:   2,
		// 	offset:  []int{0, 4},
		// },
		{
			name:    "abaaba",
			text:    "abaaba",
			pattern: "baa",
			count:   1,
			offset:  []int{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewFMIndex(tt.text)
			got := s.Search(tt.pattern)
			if !reflect.DeepEqual(got.Count(), tt.count) {
				t.Errorf("FMIndexRoweRsult.String() = %v, want %v", got.Count(), tt.count)
			}
			for _, r := range got.Rows() {
				str := r.String()
				if !reflect.DeepEqual(str, tt.pattern) {
					t.Errorf("FMIndexRoweRsult.String() = %v, want %v", str, tt.pattern)
				}

				// offset := r.Offset()
				// if !reflect.DeepEqual(offset, tt.offset[i]) {
				// 	t.Errorf("FMIndexRoweRsult.Offset() = %v, want %v", offset, tt.offset[i])
				// }
			}
		})
	}
}
