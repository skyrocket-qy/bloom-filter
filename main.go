package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
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
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	capacities := []int{10, 100, 1000, 10000, 100000, 1000000}
	realAmounts := []int{10, 100, 1000, 10000, 100000, 1000000, 10000000}
	errRates := []float64{0.0001, 0.001, 0.01, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5}

	for _, n := range capacities {
		for _, realAmount := range realAmounts {
			for _, p := range errRates {
				// calculate m and k
				m, k := BloomFilterParams(n, p)
				fmt.Printf("For n=%d items and p=%.4f:\n", n, p)
				fmt.Printf("  Bit array size (m) = %d bits (~%.2f MB)\n", m, float64(m)/8/1024/1024)
				fmt.Printf("  Number of hash functions (k) = %d\n", k)

				filterName := fmt.Sprintf("filter_n%d_r%d_p%.4f", n, realAmount, p)

				// Reserve filter
				rdb.Do(ctx, "BF.RESERVE", filterName, p, n)

				// Insert items
				startInsert := time.Now()
				for i := 0; i < realAmount && i < n; i++ {
					rdb.Do(ctx, "BF.ADD", filterName, fmt.Sprintf("item%d", i))
				}
				insertTime := time.Since(startInsert)
				fmt.Printf("  Inserted %d items in %v\n", realAmount, insertTime)

				// Test membership for existing and non-existing items
				hits, falsePositives := 0, 0
				startCheck := time.Now()
				for i := 0; i < realAmount; i++ {
					exists, _ := rdb.Do(ctx, "BF.EXISTS", filterName, fmt.Sprintf("item%d", i)).Int()
					if exists == 1 {
						hits++
					}
				}

				// Test some non-existent items for false positive rate
				testSize := 1000
				for i := 0; i < testSize; i++ {
					exists, _ := rdb.Do(ctx, "BF.EXISTS", filterName, fmt.Sprintf("fakeitem%d", i)).Int()
					if exists == 1 {
						falsePositives++
					}
				}
				checkTime := time.Since(startCheck)

				fmt.Printf("  Hit rate: %d/%d, False positives: %d/%d\n", hits, realAmount, falsePositives, testSize)
				fmt.Printf("  Query time: %v\n", checkTime)

				// Clean up: delete filter
				rdb.Del(ctx, filterName)
				fmt.Println("  Filter deleted")
			}
		}
	}
}
