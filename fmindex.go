package fmindex

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rossmerr/bwt"
	"github.com/rossmerr/wavelettree"
	"github.com/rossmerr/wavelettree/prefixtree"
)

type FMIndex struct {
	// first column of the BWT matrix
	f *wavelettree.WaveletTree
	// last column of the BWT matrix
	l               *wavelettree.WaveletTree
	prefix          *prefixtree.Prefix
	caseinsensitive bool
}

type FMIndexOption func(f *FMIndex)

func WithCaseInsensitive(caseinsensitive bool) FMIndexOption {
	return func(f *FMIndex) {
		f.caseinsensitive = caseinsensitive
	}
}

func WithPrefixTree(prefix prefixtree.Prefix) FMIndexOption {
	return func(f *FMIndex) {
		f.prefix = &prefix
	}
}

func NewFMIndex(text string, opts ...FMIndexOption) (*FMIndex, error) {
	index := &FMIndex{}

	for _, opt := range opts {
		opt(index)
	}

	if index.prefix == nil {
		en := prefixtree.English()
		index.prefix = &en
	}

	if index.caseinsensitive {
		text = strings.ToUpper(text)
	}

	first, last, err := bwt.BwtFirstLast(text)
	if err != nil {
		return nil, err
	}

	fmt.Println(first)
	index.f = wavelettree.NewWaveletTree(first, *index.prefix)
	index.l = wavelettree.NewWaveletTree(last, *index.prefix)

	return index, nil
}

func (s *FMIndex) Extract(offset, length int) string {
	return ""

}

func (s *FMIndex) Count(pattern string) int {
	return 0
}

func (s *FMIndex) Locate(pattern string) int {
	return 0

}

// LF Mapping, look at the last column of the BWT matrix and the follow it to the first column
func (s *FMIndex) lf(c rune, index int) int {
	rank := s.l.Rank(c, index)
	return s.f.Select(c, rank)
}

func (s *FMIndex) Search(pattern string) *FMIndexResult {
	if s.caseinsensitive {
		pattern = strings.ToUpper(pattern)
	}

	occurrences := s.firstOccurrences(pattern)

	length := len(pattern)
	count := 0

	for p := length - 1; p > 0; p-- {
		next := rune(pattern[p-1])

		fmt.Println("current " + string(next))
		toDelete := []int{}
		for i, arr := range occurrences {
			index := arr[count]

			if s.l.Access(index) == next {
				occurrences[i] = append(arr, s.lf(next, index))
			} else {
				toDelete = append(toDelete, i)
			}

		}
		sort.Ints(toDelete)

		for i := len(toDelete) - 1; i >= 0; i-- {
			o := toDelete[i]
			occurrences = append(occurrences[:o], occurrences[o+1:]...)
		}

		count++
	}

	result := &FMIndexResult{
		pattern: pattern,
		first:   rune(pattern[len(pattern)-1]),
		index:   s,
		count:   len(occurrences),
	}

	rows := []*FMIndexRoweRsult{}
	for _, arr := range occurrences {
		rows = append(rows, &FMIndexRoweRsult{
			access: arr,
			result: result,
		})
	}

	result.rows = rows

	return result
}

func (s *FMIndex) firstOccurrences(pattern string) [][]int {
	result := [][]int{}
	length := len(pattern)

	// look at the pattern in reverse order
	last := rune(pattern[length-1])

	matching := false

	// skip rows in the first BWT until you reach the last rune from the pattern
	start := s.f.Select(last, 0)
	for i := start; i < s.f.Length(); i++ {
		next := s.f.Access(i)
		if next == last {
			result = append(result, []int{i})
			matching = true
		} else if (next != last) && matching {
			break
		}
	}

	return result
}

type FMIndexResult struct {
	pattern string
	first   rune
	index   *FMIndex
	count   int
	rows    []*FMIndexRoweRsult
}

func (s *FMIndexResult) Count() int {
	return s.count
}

func (s *FMIndexResult) Rows() []*FMIndexRoweRsult {
	return s.rows
}

type FMIndexRoweRsult struct {
	result *FMIndexResult
	access []int
}

// func (s *FMIndexRoweRsult) Offset() int {
// 	first := s.access[len(s.access)-1]

// 	r := s.result.index.first.Rank(s.result.first, first)

// 	fmt.Println(r)

// 	sel := s.result.index.first.Select(s.result.first, r)
// 	fmt.Println(sel)

// 	return first
// }

func (s *FMIndexRoweRsult) String() string {
	str := ""
	for i := len(s.access) - 1; i >= 0; i-- {
		str += string(s.result.index.f.Access(s.access[i]))
	}

	return str
}
