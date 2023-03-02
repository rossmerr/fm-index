package fmindex

import (
	"fmt"
	"strings"

	"github.com/rossmerr/bwt"
	"github.com/rossmerr/bwt/suffixarray"
	"github.com/rossmerr/wavelettree"
	"github.com/rossmerr/wavelettree/prefixtree"
)

type FMIndex struct {
	// first column of the BWT matrix
	f *wavelettree.WaveletTree
	// last column of the BWT matrix
	l *wavelettree.WaveletTree
	// suffix array
	suffix suffixarray.Suffix
	// prefix tree
	prefix          *prefixtree.Prefix
	caseinsensitive bool
	compression     int
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

func WithCompression(compression int) FMIndexOption {
	return func(s *FMIndex) {
		s.compression = compression
	}
}

func NewFMIndex(text string, opts ...FMIndexOption) (*FMIndex, error) {
	index := &FMIndex{}

	for _, opt := range opts {
		opt(index)
	}

	if index.caseinsensitive {
		text = strings.ToUpper(text)
	}

	first, last, sa, err := bwt.BwtFirstLastSuffix[suffixarray.SuffixArray](text, suffixarray.WithCompression(index.compression))
	if err != nil {
		return nil, err
	}

	if index.prefix == nil {
		pt := prefixtree.NewHuffmanCodeTree(first)
		index.prefix = &pt
	}

	index.suffix = sa
	index.f = wavelettree.NewWaveletTree(first, *index.prefix)
	index.l = wavelettree.NewWaveletTree(last, *index.prefix)

	return index, nil
}

func (s *FMIndex) Extract(offset, length int) string {
	result := make([]rune, length)
	iterator := s.suffix.Enumerate()
	find := map[int]int{}
	for iterator.HasNext() {
		currentElement, index := iterator.Next()
		find[index] = currentElement
	}

	rs := NewReveriseSuffix(s.suffix)

	count := 0
	for i := offset; i < offset+length; i++ {
		if rs.Has(i) {
			r := s.f.Access(i)
			//s.prefix[]
			result[count] = r
			count++
		}
		p := rs.Walk(i)
		//	s.f.
		//l := find[i]
		r := s.f.Access(p)
		s.f.Rank(r, p)

		//s.l.
		fmt.Println((string(r)))
		result[count] = r
		count++
	}
	return string(result)
	//return ""

}

func (s *FMIndex) Count(pattern string) int {
	f, l := s.query(pattern)
	return l - f
}

func (s *FMIndex) Locate(pattern string) []int {
	f, l := s.query(pattern)
	result := []int{}
	for i := f; i < l; i++ {
		r := s.suffix.Get(i)
		result = append(result, r)
	}
	if f == l {
		r := s.suffix.Get(f)
		result = append(result, r)
	}
	return result
}

func (s *FMIndex) query(pattern string) (top, bottom int) {
	if s.caseinsensitive {
		pattern = strings.ToUpper(pattern)
	}

	length := len(pattern)

	// // look at the pattern in reverse order
	next := rune(pattern[length-1])

	n1, _ := s.f.Rank(next, 0)
	top = s.f.Select(next, n1)

	n2, _ := s.f.Rank(next, s.l.Length())
	bottom = s.f.Select(next, n2+1)

	i := length - 2
	for i >= 0 && bottom >= top {
		next = rune(pattern[i])
		n1, _ := s.l.Rank(next, top)
		n2, _ := s.l.Rank(next, bottom)
		skip := s.f.Select(next, 0)
		top = (n1 + skip)
		bottom = (n2 + skip)
		i--
	}

	return
}

type ReveriseSuffix struct {
	find map[int]int
}

func NewReveriseSuffix(suffix suffixarray.Suffix) *ReveriseSuffix {
	iterator := suffix.Enumerate()
	find := map[int]int{}
	for iterator.HasNext() {
		currentElement, index := iterator.Next()
		find[index] = currentElement
	}
	return &ReveriseSuffix{
		find: find,
	}
}

func (s *ReveriseSuffix) Walk(index int) int {
	for _, ok := s.find[index]; ok; index-- {
		return index
	}

	return -1
}

func (s *ReveriseSuffix) Has(index int) bool {
	_, ok := s.find[index]
	return ok
}
