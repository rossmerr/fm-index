package fmindex

import (
	"strings"

	"github.com/rossmerr/bwt"
	"github.com/rossmerr/wavelettree"
	"github.com/rossmerr/wavelettree/prefixtree"
)

type FMIndex struct {
	first           *wavelettree.WaveletTree
	last            *wavelettree.WaveletTree
	caseinsensitive bool
}

func NewFMIndex(text string, caseinsensitive bool) (*FMIndex, error) {
	if caseinsensitive {
		text = strings.ToUpper(text)
	}
	first, last, err := bwt.BwtFirstLast(text)
	if err != nil {
		return nil, err
	}

	en := prefixtree.English()
	f := wavelettree.NewWaveletTree(first, en)
	l := wavelettree.NewWaveletTree(last, en)

	return &FMIndex{
		first:           f,
		last:            l,
		caseinsensitive: caseinsensitive,
	}, nil
}

func (s *FMIndex) Search(pattern string) *FMIndexResult {
	if s.caseinsensitive {
		pattern = strings.ToUpper(pattern)
	}

	occurrences := s.firstOccurrences(pattern)

	length := len(pattern) - 1

	for p := length - 1; p > 0; p-- {
		c := rune(pattern[p])
		n := rune(pattern[p-1])

		arr := occurrences[p]
		for i, rank := range arr {
			r := s.first.Select(c, rank)
			a := s.last.Access(r)
			if a == n {
				rank := s.last.Rank(a, r)
				// results for the next match
				occurrences[p-1] = append(occurrences[p-1], rank)
			} else {
				arr = append(arr[:i], arr[i+1:]...)
			}
		}
	}

	result := &FMIndexResult{}
	rows := []*FMIndexResultRow{}
	if _, ok := occurrences[0]; ok {
		for _, arr := range occurrences {
			rows = append(rows, &FMIndexResultRow{
				start: arr[0],
				end:   arr[length-1],
			})
		}
	}

	result.count = len(rows)
	result.rows = rows

	return result
}

func (s *FMIndex) firstOccurrences(pattern string) map[int][]int {
	result := map[int][]int{}
	length := len(pattern) - 1

	// look at the pattern in reverse order
	last := rune(pattern[length])

	matching := false
	// look over each item in Wavelet Tree until it matches the last rune
	// todo can rank be used from the Wavelet Tree to do this?
	for index := 0; index < s.first.Length(); index++ {
		i := s.first.Access(index)
		if i == last {
			rank := s.first.Rank(last, index)
			result[length-1] = append(result[length-1], rank)
			matching = true
		} else if (i != last) && matching {
			break
		}
		index++
	}

	return result
}

type FMIndexResult struct {
	count int
	rows  []*FMIndexResultRow
}

func (s *FMIndexResult) Count() int {
	return s.count
}

func (s *FMIndexResult) Rows() []*FMIndexResultRow {
	return s.rows
}

type FMIndexResultRow struct {
	start int
	end   int
}

func (s *FMIndexResultRow) Start() int {
	return s.start
}

func (s *FMIndexResultRow) End() int {
	return s.end
}
