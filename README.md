# FM-index

[![Go](https://github.com/rossmerr/fm-index/actions/workflows/go.yml/badge.svg)](https://github.com/rossmerr/fm-index/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rossmerr/fm-index)](https://goreportcard.com/report/github.com/rossmerr/fm-index)
[![Read the Docs](https://pkg.go.dev/badge/golang.org/x/pkgsite)](https://pkg.go.dev/github.com/rossmerr/fm-index)

A [FM-index](https://en.wikipedia.org/wiki/FM-index) using [Wavelet Tree's](https://en.wikipedia.org/wiki/Wavelet_Tree) to store the [Burrows–Wheeler transform](https://en.wikipedia.org/wiki/Burrows%E2%80%93Wheeler_transform). A [Prefix Tree](https://en.wikipedia.org/wiki/Trie) using [BitVector](https://en.wikipedia.org/wiki/Bit_array) and a Sampled Suffix Array, for storing the offsets.

The Prefix Tree can be shared across FM-indexs to reduce storage needs.

```go
index, err := NewFMIndex("The quick brown fox jumps over the lazy dog", WithCompression(2))
```

## Operation

### Count

```go
count := index.Count("jumps")
fmt.Println(count) // 1
```

### Locate

```go
matches := index.Locate("jumps")
fmt.Println(matches) // []int{20}
```

### Extract

```go
text := index.Extract(matches[0], 5)
fmt.Println(text) // "jumps"
```
