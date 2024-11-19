package util

import (
	"bytes"
	"fmt"
	"sort"
)

// BoxWhisker : a Box Whisker plot
type BoxWhisker struct {
	// data points
	data []int
	// number of data point samples
	n int
	// box whisker statistics
	min, max, median, q1, q3 int
}

// NewBoxWhisker : create a Box Whisker of a given data size
func NewBoxWhisker(size int) (*BoxWhisker, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid size %d", size)
	}
	return &BoxWhisker{
		data: make([]int, size),
		n:    0,
	}, nil
}

// GetNumSamples : get the number of smples
func (bw *BoxWhisker) GetNumSamples() int {
	return bw.n
}

// AddSample : add a data point
func (bw *BoxWhisker) AddSample(x int) {
	if bw.n < len(bw.data) {
		bw.data[bw.n] = x
		bw.n++
	}
}

// Calculate : calculate statistics
func (bw *BoxWhisker) Calculate() {
	if bw.n == 0 {
		bw.min = 0
		bw.max = 0
		bw.median = 0
		bw.q1 = 0
		bw.q3 = 0
	}

	sorted := bw.data[0:bw.n]
	sort.Ints(sorted)

	bw.min = sorted[0]
	bw.max = sorted[bw.n-1]
	bw.median = sorted[bw.n/2]
	bw.q1 = sorted[bw.n/4]
	bw.q3 = sorted[3*bw.n/4]
}

// String : a print out of the Box Whisker statistics
func (bw *BoxWhisker) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "numSamples=%d; ", bw.n)
	fmt.Fprintf(&b, "median=%v; ", bw.median)
	fmt.Fprintf(&b, "min=%v; ", bw.min)
	fmt.Fprintf(&b, "max=%v; ", bw.max)
	fmt.Fprintf(&b, "q1=%v; ", bw.q1)
	fmt.Fprintf(&b, "q3=%v; ", bw.q3)
	return b.String()
}
