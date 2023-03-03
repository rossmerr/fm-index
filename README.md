# FM-index

[![Go](https://github.com/rossmerr/fm-index/actions/workflows/go.yml/badge.svg)](https://github.com/rossmerr/fm-index/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rossmerr/fm-index)](https://goreportcard.com/report/github.com/rossmerr/fm-index)
[![Read the Docs](https://pkg.go.dev/badge/golang.org/x/pkgsite)](https://pkg.go.dev/github.com/rossmerr/fm-index)

A [FM-index](https://en.wikipedia.org/wiki/FM-index) using [Wavelet Tree's](https://en.wikipedia.org/wiki/Wavelet_Tree) to store the [Burrows–Wheeler transform](https://en.wikipedia.org/wiki/Burrows%E2%80%93Wheeler_transform). A [Prefix Tree](https://en.wikipedia.org/wiki/Trie) using [BitVector](https://en.wikipedia.org/wiki/Bit_array) and a Sampled Suffix Array, for storing the offsets.

The Prefix Tree can be shared across FM-indexs to reduce storage needs.

Implements the following operation

- Count
- Locate
- Extract
