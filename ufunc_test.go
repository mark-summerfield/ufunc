package ufunc

import (
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/mark-summerfield/unum"
)

func Test_LastIndex(t *testing.T) {
	//         0  1  2  3   4   5   6  7  8  9
	a := []int{2, 4, 6, 8, 10, 12, 10, 8, 6, 4}
	e := 3
	if i := slices.Index(a, 8); i != e {
		t.Errorf("expected %d; got %d", e, i)
	}
	e = 7
	if i := LastIndex(a, 8); i != e {
		t.Errorf("expected %d; got %d", e, i)
	}
	e = -1
	if i := LastIndex(a, 88); i != e {
		t.Errorf("expected %d; got %d", e, i)
	}
}

func Test_LastIndexFunc(t *testing.T) {
	//         0  1  2  3   4   5   6  7  8  9
	a := []int{2, 4, 6, 8, 10, 12, 10, 8, 6, 4}
	e := 7
	if i := LastIndexFunc(a, func(x int) bool { return x == 8 }); i != e {
		t.Errorf("expected %d; got %d", e, i)
	}
	e = -1
	if i := LastIndexFunc(a, func(x int) bool { return x == 88 }); i != e {
		t.Errorf("expected %d; got %d", e, i)
	}
}

func Test_HasPrefix(t *testing.T) {
	//         0  1  2  3   4   5   6  7  8  9
	a := []int{2, 4, 6, 8, 10, 12, 10, 8, 6, 4}
	p := []int{2, 4, 6}
	if !HasPrefix(a, p) {
		t.Error("expected true; got false")
	}
	p[2] = 7
	if HasPrefix(a, p) {
		t.Error("expected false; got true")
	}
}

func Test_HasSuffix(t *testing.T) {
	//         0  1  2  3   4   5   6  7  8  9
	a := []int{2, 4, 6, 8, 10, 12, 10, 8, 6, 4}
	s := []int{8, 6, 4}
	if !HasSuffix(a, s) {
		t.Error("expected true; got false")
	}
	s[0] = 7
	if HasSuffix(a, s) {
		t.Error("expected false; got true")
	}
}

func Test_Map(t *testing.T) {
	// tag::mapxeg1[]
	reals := []float64{1.2, -4, 8.5, 19.6, 14.2, -15.5, 18.7}
	// end::mapxeg1[]
	posInts := make([]int, 0, len(reals))
	processPosInt := func(i int) { posInts = append(posInts, i) }
	// tag::mapxeg[]
	for i := range Map(reals, func(x float64) (int, bool) {
		if x < 0 {
			return 0, false
		}
		return int(math.Round(x)), true
	}) {
		processPosInt(i)
	}
	// end::mapxeg[]
	expected := []int{1, 9, 20, 14, 19}
	if slices.Compare(expected, posInts) != 0 {
		t.Errorf("expected %v, got %v", expected, posInts)
	}
}

func Test_manual_span(t *testing.T) {
	ints := []int{1, 4, 9, 16, 25, 36, 49, 64, 81, 100, 121}
	subs := [][]int{}
	for i := 0; i < len(ints); i += 4 {
		subs = append(subs, ints[i:min(i+4, len(ints))])
	}
	exp := "[[1 4 9 16] [25 36 49 64] [81 100 121]]"
	got := fmt.Sprintf("%v", subs)
	if exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	subs = subs[:0]
	for i := 0; i < len(ints); i += 3 {
		subs = append(subs, ints[i:min(i+3, len(ints))])
	}
	exp = "[[1 4 9] [16 25 36] [49 64 81] [100 121]]"
	got = fmt.Sprintf("%v", subs)
	if exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	subs = subs[:0]
	for i := 0; i < 10; i += 2 {
		subs = append(subs, ints[i:min(i+2, len(ints))])
	}
	exp = "[[1 4] [9 16] [25 36] [49 64] [81 100]]"
	got = fmt.Sprintf("%v", subs)
	if exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func Test_Span(t *testing.T) {
	// tag::span0[]
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	// end::span0[]
	spans := [][]int{}
	lastok := true
	// tag::span1[]
	for span, ok := range Spans(data, 4) { // <1>
		// [1 2 3 4]:true [5 6 7 8]:true [9 10 11]:false
		// end::span1[]
		spans = append(spans, span)
		lastok = ok
	}
	if lastok {
		t.Error("expeced last ok false")
	}
	exp := "[[1 2 3 4] [5 6 7 8] [9 10 11]]"
	got := fmt.Sprintf("%v", spans)
	if exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	spans = spans[:0]
	// tag::span2[]
	for span, ok := range Spans(data, 3) { // <2>
		// [1 2 3]:true [4 5 6]:true [7 8 9]:true [10 11]:false
		// end::span2[]
		spans = append(spans, span)
		lastok = ok
	}
	if lastok {
		t.Error("expeced last ok false")
	}
	exp = "[[1 2 3] [4 5 6] [7 8 9] [10 11]]"
	got = fmt.Sprintf("%v", spans)
	if exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	spans = spans[:0]
	// tag::span3[]
	for span, ok := range Spans(data[:10], 2) { // <3>
		// [1 2]:true [3 4]:true [5 6]:true [7 8]:true [9 10]:true
		// end::span3[]
		spans = append(spans, span)
		lastok = ok
	}
	if !lastok {
		t.Error("expeced last ok true")
	}
	exp = "[[1 2] [3 4] [5 6] [7 8] [9 10]]"
	got = fmt.Sprintf("%v", spans)
	if exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func Test_Range_int(t *testing.T) {
	ints := make([]int, 0, 10)
	// tag::range_int[]
	for i := range 10 { // 0 1 2 … 9 <2>
		// end::range_int[]
		ints = append(ints, i)
	}
	ix := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if slices.Compare(ix, ints) != 0 {
		t.Errorf("expected %v; got %v", ix, ints)
	}
	ints = ints[:0] // ints: []int{}
	// tag::all2[]
	for i := range Range(5, 15) { // 5 6 7 … 14
		// end::all2[]
		ints = append(ints, i)
	}
	ix = []int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	if slices.Compare(ix, ints) != 0 {
		t.Errorf("expected %v; got %v", ix, ints)
	}
	ints = ints[:0] // ints: []int{}
	total := 0.0
	// tag::all3[]
	for n := range Range(40.0, 81.0) { // 40.0 41.0 42.0 … 80.0
		// end::all3[]
		total += n
	}
	ex := 2460.0
	if !unum.IsClose(ex, total) {
		t.Errorf("expected %f; got %f", ex, total)
	}
	ints = ints[:0] // ints: []int{}
	// tag::all4[]
	for i := range RangeX(10, 32, 3) { // 10 13 16 … 31
		// end::all4[]
		ints = append(ints, i)
	}
	ix = []int{10, 13, 16, 19, 22, 25, 28, 31}
	if slices.Compare(ix, ints) != 0 {
		t.Errorf("expected %v; got %v", ix, ints)
	}
	ints = ints[:0] // ints: []int{}
	for i := range RangeX(30, 9, 3) {
		ints = append(ints, i)
	}
	ix = []int{30, 27, 24, 21, 18, 15, 12}
	if slices.Compare(ix, ints) != 0 {
		t.Errorf("expected %v; got %v", ix, ints)
	}
	ints = ints[:0] // ints: []int{}
	for i := range RangeX(27, 10, 4) {
		ints = append(ints, i)
	}
	ix = []int{27, 23, 19, 15, 11}
	if slices.Compare(ix, ints) != 0 {
		t.Errorf("expected %v; got %v", ix, ints)
	}
	ints = ints[:0] // ints: []int{}
	for i := range Range(0, 0) {
		ints = append(ints, i)
	}
	ix = []int{}
	if slices.Compare(ix, ints) != 0 {
		t.Errorf("expected %v; got %v", ix, ints)
	}
	runes := []rune{}
	rx := []rune{'M', 'N', 'O', 'P', 'Q'}
	// tag::all7[]
	for r := range Range('M', 'R') { // 'M' 'N' 'O' 'P' 'Q' <2>
		// end::all7[]
		runes = append(runes, r)
	}
	if slices.Compare(rx, runes) != 0 {
		t.Errorf("expected %v; got %v", rx, runes)
	}
}

func Test_Range_real(t *testing.T) {
	reals := make([]float64, 0, 10)
	for i := range Range(5.0, 15.0) {
		reals = append(reals, i)
	}
	ix := []float64{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	if !equalReals(ix, reals) {
		t.Errorf("expected %v; got %v", ix, reals)
	}
	reals = reals[:0] // reals: []float64{}
	for i := range Range(7.0, 15.0) {
		reals = append(reals, i)
	}
	ix = []float64{7, 8, 9, 10, 11, 12, 13, 14}
	if !equalReals(ix, reals) {
		t.Errorf("expected %v; got %v", ix, reals)
	}
	reals = reals[:0] // reals: []float64{}
	// tag::all5[]
	for i := range RangeX(10.0, 15.0, 0.5) { // 10.0 10.5 11.0 … 14.5
		// end::all5[]
		reals = append(reals, i)
	}
	ix = []float64{10, 10.5, 11, 11.5, 12, 12.5, 13, 13.5, 14, 14.5}
	if !equalReals(ix, reals) {
		t.Errorf("expected %v; got %v", ix, reals)
	}
	reals = reals[:0] // reals: []float64{}
	// tag::all6[]
	for n := range RangeX(20.0, 9.5, 0.5) { // 20.0 19.5 19.0 … 10.0 <1>
		// end::all6[]
		reals = append(reals, n)
	}
	ix = []float64{
		20, 19.5, 19, 18.5, 18, 17.5, 17, 16.5, 16, 15.5, 15,
		14.5, 14, 13.5, 13, 12.5, 12, 11.5, 11, 10.5, 10,
	}
	if !equalReals(ix, reals) {
		t.Errorf("expected %v; got %v", ix, reals)
	}
	reals = reals[:0] // reals: []float64{}
	for i := range RangeX(27.0, 22.0, 1.5) {
		reals = append(reals, i)
	}
	ix = []float64{27, 25.5, 24, 22.5}
	if !equalReals(ix, reals) {
		t.Errorf("expected %v; got %v", ix, reals)
	}
	reals = reals[:0] // reals: []float64{}
	for i := range RangeX(1.0, 8.5, 0.5) {
		reals = append(reals, i)
	}
	ix = []float64{
		1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 5.5, 6, 6.5,
		7, 7.5, 8,
	}
	if !equalReals(ix, reals) {
		t.Errorf("expected %v; got %v", ix, reals)
	}
}

func equalReals(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if !unum.IsClose(x, b[i]) {
			return false
		}
	}
	return true
}

func Test_Merge(t *testing.T) {
	var ns []int
	// tag::mergeeg[]
	for n := range Merge(Range(0, 10), Range(10, 20), Range(20, 30)) {
		// end::mergeeg[]
		ns = append(ns, n)
	}
	exp := []int{
		0, 10, 20, 1, 11, 21, 2, 12, 22, 3, 13, 23, 4, 14, 24, 5,
		15, 25, 6, 16, 26, 7, 17, 27, 8, 18, 28, 9, 19, 29,
	}
	if slices.Compare(ns, exp) != 0 {
		t.Errorf("expected %v, got %v", exp, ns)
	}
}

func Test_Zip(t *testing.T) {
	var ns [][]int
	// tag::zipeg[]
	for row := range Zip(RangeX(0, 11, 3), RangeX(1, 11, 3),
		RangeX(2, 11, 3)) {
		// [0 1 2] [3 4 5] [6 7 8]
		// end::zipeg[]
		ns = append(ns, row)
	}
	exp := [][]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}}
	for i, en := range exp {
		if slices.Compare(ns[i], en) != 0 {
			t.Errorf("expected %v, got %v", en, ns[i])
		}
	}
}

func Test_ZipLongest(t *testing.T) {
	var ns [][]int
	// tag::ziplongeg[]
	for row := range ZipLongest(RangeX(0, 11, 3), RangeX(1, 11, 3),
		RangeX(2, 11, 3)) {
		// [0 1 2 0] [3 4 5 0] [6 7 8 9]
		// end::ziplongeg[]
		ns = append(ns, row)
	}
	exp := [][]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {9, 10, 0}}
	for i, en := range exp {
		if slices.Compare(ns[i], en) != 0 {
			t.Errorf("expected %v, got %v", en, ns[i])
		}
	}
}
