package hyperloglog

import (
	"fmt"
	"math"
	"math/bits"

	"github.com/cespare/xxhash/v2"
)

type HyperLogLog struct {
	buckets []int // store the maximum
	b       uint
}

var CONSTANT = 0.72134

func New(b uint) *HyperLogLog {
	// guide for what b to use, m = 2^b, accuracy = 1.04 / sqrt(m)
	m := math.Pow(2, float64(b)) // hard limit for the b bits, b=4 is 1111=15 with m = 16
	return &HyperLogLog{buckets: make([]int, int(m)), b: uint(b)}
}

func NewWithErrorRate(error_rate float64) *HyperLogLog {
	// log2 to get the next power of 2 above raw m
	// ceil to round up and get the above
	// math.pow 2^12 || 1 << 12
	raw_m := math.Pow((1.04 / error_rate), 2)
	m := math.Pow(2, math.Ceil(math.Log2(raw_m)))
	b := math.Log2(m)
	return New(uint(b))
}

func (h *HyperLogLog) Add(v string) {
	// hash -> get first b bits and assign to a bucket ->
	// get w remaining bits and get the leading zeros -> compare the current max with the leading zeros
	// x := hash(v)

	x := xxhash.Sum64String(v)
	fmt.Println("x", x)
	// // j:= v.substr(b) + 1; b = log base 2 of (m); b first b bits as identifier
	bucket := x >> (64 - h.b)
	// p(w) = position of leftmost 1;  w = remaining bits;
	w := bits.LeadingZeros64(x<<uint64(h.b)) + 1

	// get the max
	h.buckets[bucket] = max(h.buckets[bucket], w)
}

func (h *HyperLogLog) Count() int {
	harmonicSum := 0.0
	for _, i := range h.buckets {
		harmonicSum += 1.0 / math.Pow(2, float64(i))
	}

	m := float64(len(h.buckets))
	alpha := 0.7213 / (1 + 1.079/m)
	return int(alpha * m * (m / harmonicSum))
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
