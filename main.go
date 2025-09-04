package main

import (
	"fmt"
	"math"
)

// BloomFilterParams calculates m and k for a given n and p
func BloomFilterParams(n int, p float64) (m int, k int) {
	if n <= 0 || p <= 0 || p >= 1 {
		panic("n must be > 0 and 0 < p < 1")
	}

	// m = -(n * ln(p)) / (ln 2)^2
	mFloat := -float64(n) * math.Log(p) / (math.Ln2 * math.Ln2)
	m = int(math.Ceil(mFloat)) // round up to next integer

	// k = (m/n) * ln 2
	kFloat := (float64(m) / float64(n)) * math.Ln2
	k = int(math.Ceil(kFloat)) // round up to next integer

	return
}

func main() {
	n := 1000000 // expected items
	p := 0.01    // false positive probability

	m, k := BloomFilterParams(n, p)
	fmt.Printf("For n=%d items and p=%.4f:\n", n, p)
	fmt.Printf("  Bit array size (m) = %d bits (~%.2f MB)\n", m, float64(m)/8/1024/1024)
	fmt.Printf("  Number of hash functions (k) = %d\n", k)
}
