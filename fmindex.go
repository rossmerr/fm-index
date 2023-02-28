package fmindex

import (
	"strings"

	"github.com/rossmerr/bwt"
	"github.com/rossmerr/wavelettree"
	"github.com/rossmerr/wavelettree/prefixtree"
)

type FMIndex struct {
	// first column of the BWT matrix
	f *wavelettree.WaveletTree
	// last column of the BWT matrix
	l *wavelettree.WaveletTree
	// suffix array
	sa bwt.Suffix
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

	if index.prefix == nil {
		en := prefixtree.English()
		index.prefix = &en
	}

	if index.caseinsensitive {
		text = strings.ToUpper(text)
	}

	first, last, sa, err := bwt.BwtFirstLastSuffix[bwt.SampleSuffixArray](text, bwt.WithCompression(index.compression))
	if err != nil {
		return nil, err
	}
	index.sa = sa
	index.f = wavelettree.NewWaveletTree(first, *index.prefix)
	index.l = wavelettree.NewWaveletTree(last, *index.prefix)

	return index, nil
}

func (s *FMIndex) Extract(offset, length int) string {
	return ""

}

func (s *FMIndex) Count(pattern string) int {
	f, l := s.query(pattern)
	return l - f
}

func (s *FMIndex) Locate(pattern string) []int {
	f, l := s.query(pattern)
	result := []int{}
	for i := f; i < l; i++ {
		result = append(result, s.findSuffix(i, 0))
	}
	return result
}

func (s *FMIndex) findSuffix(i, count int) int {
	if r, ok := s.sa.Get(i); ok {
		return r + count
	} else {
		return s.findSuffix(i-1, count+1)
	}
}

func (s *FMIndex) query(pattern string) (top, bottom int) {
	if s.caseinsensitive {
		pattern = strings.ToUpper(pattern)
	}

	length := len(pattern)

	// // look at the pattern in reverse order
	next := rune(pattern[length-1])

	top = s.f.Select(next, 0)
	bottom = s.f.Select(next, s.l.Length()) + 1

	i := length - 2
	for i >= 0 && bottom > top {
		next = rune(pattern[i])
		n1 := s.l.Rank(next, top)
		n2 := s.l.Rank(next, bottom)
		skip := s.f.Select(next, 0)
		top = (n1 + skip)
		bottom = (n2 + skip)
		i--
	}

	return
}
