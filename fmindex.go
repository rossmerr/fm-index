package fmindex

import (
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
	compression     *int
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
		if compression >= 2 {
			s.compression = &compression
		}
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

	if index.compression == nil {
		two := 2
		index.compression = &two
	}

	first, last, sa, err := bwt.FirstLastSuffix(text, bwt.WithSampleRate(*index.compression))
	if err != nil {
		return nil, err
	}

	if index.prefix == nil {
		index.prefix = prefixtree.NewHuffmanCodeTree(first)
	}

	index.suffix = sa
	index.f = wavelettree.NewWaveletTree(first, index.prefix)
	index.l = wavelettree.NewWaveletTree(last, index.prefix)

	return index, nil
}

func (s *FMIndex) Extract(offset, length int) string {
	result := []rune{}

	mappedSuffix := map[int]int{}
	iterator := s.suffix.Enumerate()
	for iterator.HasNext() {
		k, i := iterator.Next()
		mappedSuffix[k] = i
	}

	index := 0
	ok := false
	for i := offset; i < offset+length; i++ {
		if index, ok = mappedSuffix[i]; !ok {
			index = i - 1
			hops := 1

			for {
				if i, ok := mappedSuffix[index]; ok {
					index = i
					break
				}
				index--
				hops++
			}

			for i := 0; i < hops; i++ {
				r := s.f.Access(index)
				rank, _ := s.f.Rank(r, index)
				index = s.l.Select(r, rank)
			}
		}

		r := s.f.Access(index)
		result = append(result, r)

	}
	return string(result)
}

func (s *FMIndex) Count(pattern string) int {
	f, l := s.query(pattern)
	return l - f
}

func (s *FMIndex) Locate(pattern string) []int {
	f, l := s.query(pattern)
	result := []int{}
	for i := f; i < l; i++ {
		index := s.walkToNearest(i, 0)
		r := s.suffix.Get(index)
		result = append(result, r)
	}
	if f == l {
		index := s.walkToNearest(f, 0)

		r := s.suffix.Get(index)
		result = append(result, r)
	}
	return result
}

func (s *FMIndex) walkToNearest(index, count int) int {
	b := s.suffix.Has(index)
	if b {
		return index + count
	}
	count++
	a := s.l.Access(index)
	r, _ := s.l.Rank(a, index)
	nextIndex := s.f.Select(a, r)
	return s.walkToNearest(nextIndex, count)
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

func (s *FMIndex) CaseInsensitive() bool {
	return s.caseinsensitive
}

func (s *FMIndex) PrefixTree() *prefixtree.Prefix {
	return s.prefix

}

func (s *FMIndex) Compression() int {
	return *s.compression

}
