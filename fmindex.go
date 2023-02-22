package fmindex

import (
	"fmt"
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
	f, l := s.query(pattern)
	return l - f
}

func (s *FMIndex) Locate(pattern string) int {
	return 0

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

		fmt.Println("current " + string(next))
		n1 := s.l.Rank(next, top)
		n2 := s.l.Rank(next, bottom)
		fmt.Println(n1)
		fmt.Println(n2)
		skip := s.f.Select(next, 0)
		top = (n1 + skip)
		bottom = (n2 + skip)
		i--
	}

	return
}
