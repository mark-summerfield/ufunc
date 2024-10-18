// Copyright © 2024 Mark Summerfield. All rights reserved.

// This package provides some generic range functions (including Map,
// Reduce, Zip, and ZipLongest) and fills in some gaps in the slices
// package. ([TOC])
//
// [TOC]: file:///home/mark/app/golib/doc/index.html
package ufunc

import (
	_ "embed"
	"iter"
	"slices"

	"github.com/mark-summerfield/ureal"
)

//go:embed Version.dat
var Version string

// HasPrefix returns true if the values slice starts with prefix
func HasPrefix[E comparable](values, prefix []E) bool {
	for i := range len(prefix) {
		if i == len(values) || values[i] != prefix[i] {
			return false
		}
	}
	return true
}

// HasSuffix returns true if the values slice ends with prefix
func HasSuffix[E comparable](values, suffix []E) bool {
	if len(suffix) > len(values) {
		return false
	}
	j := len(values) - 1
	for i := len(suffix) - 1; i >= 0; i-- {
		if values[j] != suffix[i] {
			return false
		}
		j--
	}
	return true
}

// LastIndex returns the index position of the rightmost value in the slice
// or -1 if value isn't in the slice.
// See also [slices.Index]
func LastIndex[E comparable](values []E, value E) int {
	for i := len(values) - 1; i >= 0; i-- {
		if values[i] == value {
			return i
		}
	}
	return -1
}

// LastIndex returns the index position of the rightmost value in the slice
// or -1 if value isn't in the slice.
// See also [slices.IndexFunc]
func LastIndexFunc[E any](values []E, found func(E) bool) int {
	for i := len(values) - 1; i >= 0; i-- {
		if found(values[i]) {
			return i
		}
	}
	return -1
}

// Map returns an iterator which yields every value in the sources
// transformed by the mapper function (but dropping any values for which the
// mapper's ok is false).
func Map[S, T any](sources []S, mapper func(S) (T, bool)) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, source := range sources {
			if target, ok := mapper(source); ok {
				if !yield(target) {
					return
				}
			}
		}
	}
}

// Merge accepts any number of iterators (rangefuncs) and returns a single
// iterator that yields all the first iterator's elements, then all the
// second iterator's, and so on.
func Merge[E any](rfns ...iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		pulls := make([]func() (E, bool), 0, len(rfns))
		oks := make([]bool, 0, len(rfns))
		for _, rfn := range rfns {
			pull, stop := iter.Pull(rfn)
			defer stop()
			pulls = append(pulls, pull)
			oks = append(oks, true)
		}
		for slices.Contains(oks, true) {
			for i, pull := range pulls {
				if oks[i] {
					if element, ok := pull(); ok {
						if !yield(element) {
							return
						}
					} else {
						oks[i] = false
					}
				}
			}
		}
	}
}

// Range is a range function that returns a function that
// returns numbers from start upto (or downto) the step
// before end in steps of 1.
//
//	for x := range Range(5, 15) { // 5 6 7 … 14
func Range[N ureal.SignedNumber](start, end N) iter.Seq[N] {
	return RangeX(start, end, 1)
}

// RangeX is a rang function that returns a function that
// returns numbers from start upto (or downto) the step
// before end in steps of step.
// The step must be a magnitude > 0 or RangeX will panic.
//
//	for x := range RangeX(1.0, 8.5, 0.5) { // 1.0 1.5 2.0 … 8.0
func RangeX[N ureal.SignedNumber](start, end, step N) iter.Seq[N] {
	if step <= 0 {
		panic("step size must be > 0")
	}
	return func(yield func(N) bool) {
		if start <= end {
			for ; start < end; start += step {
				if !yield(start) {
					return
				}
			}
		} else {
			for ; start > end; start -= step {
				if !yield(start) {
					return
				}
			}
		}
	}
}

// Reduce returns the accumulated elements based on the reduce function and
// the initial accumulator value.
func Reduce[E, A any](elements []E, reduce func(E, A) A, accumulator A) A {
	for _, element := range elements {
		accumulator = reduce(element, accumulator)
	}
	return accumulator
}

// Returns subslices of stride size from the given slice. Each subslice is
// returned with true, except possibly the last which if short is returned
// with false.
func Spans[T any](slice []T, stride int) iter.Seq2[[]T, bool] {
	if stride <= 0 {
		panic("stride must be > 0")
	}
	return func(yield func([]T, bool) bool) {
		for i := 0; i < len(slice); i += stride {
			subslice := slice[i:min(i+stride, len(slice))]
			if !yield(subslice, len(subslice) == stride) {
				return
			}
		}
	}
}

// Zip accepts any number of iterators (rangefuncs) and returns a single
// iterator that returns a slice of all the first elements from all the
// iterators, then a slice of all the second elements, and so on, stopping
// as soon as one of the iterators runs out.
func Zip[E any](rfns ...iter.Seq[E]) iter.Seq[[]E] {
	return func(yield func([]E) bool) {
		pulls := make([]func() (E, bool), 0, len(rfns))
		for _, rfn := range rfns {
			pull, stop := iter.Pull(rfn)
			defer stop()
			pulls = append(pulls, pull)
		}
		for {
			row := make([]E, 0, len(pulls))
			for _, pull := range pulls {
				if element, ok := pull(); ok {
					row = append(row, element)
				} else {
					return // finish when first of rfns is done
				}
			}
			if !yield(row) { // yield complete rows only
				return
			}
		}
	}
}

// ZipLongest accepts any number of iterators (rangefuncs) and returns a
// single iterator that returns a slice of all the first elements from all
// the iterators, then a slice of all the second elements, and so on. This
// continues so long as at least one iterator has values, with exhausted
// iterators' elements being replaced with their zero values.
func ZipLongest[E any](rfns ...iter.Seq[E]) iter.Seq[[]E] {
	return func(yield func([]E) bool) {
		pulls := make([]func() (E, bool), 0, len(rfns))
		for _, rfn := range rfns {
			pull, stop := iter.Pull(rfn)
			defer stop()
			pulls = append(pulls, pull)
		}
		var zero E
		for {
			row := make([]E, 0, len(pulls))
			oks := 0
			for _, pull := range pulls {
				if element, ok := pull(); ok {
					oks++
					row = append(row, element)
				} else {
					row = append(row, zero)

				}
			}
			if oks == 0 || !yield(row) {
				return
			}
		}
	}
}
