package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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

// testConfig holds the configuration for a single Bloom filter test run.
type testConfig struct {
	capacity   int
	errorRate  float64
	realAmount int
}

// testBloomFilter runs a single Bloom filter test with the given configuration.
func testBloomFilter(ctx context.Context, rdb *redis.Client, config testConfig) error {
	const testSize = 1000

	// Calculate Bloom filter parameters.
	m, k := BloomFilterParams(config.capacity, config.errorRate)
	fmt.Printf("Testing n=%d, p=%.4f, realAmount=%d -> m=%d, k=%d\n",
		config.capacity, config.errorRate, config.realAmount, m, k)

	filterName := fmt.Sprintf("filter_n%d_r%d_p%.4f", config.capacity, config.realAmount, config.errorRate)
	defer rdb.Del(ctx, filterName) // Ensure cleanup

	// Reserve the filter.
	if err := rdb.Do(ctx, "BF.RESERVE", filterName, config.errorRate, config.capacity).Err(); err != nil {
		return fmt.Errorf("failed to reserve bloom filter: %w", err)
	}

	// Insert items using a pipeline for efficiency.
	pipe := rdb.Pipeline()
	insertCount := min(config.capacity, config.realAmount)
	startInsert := time.Now()
	for i := 0; i < insertCount; i++ {
		pipe.Do(ctx, "BF.ADD", filterName, fmt.Sprintf("item%d", i))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert items: %w", err)
	}
	insertTime := time.Since(startInsert)
	fmt.Printf("  Inserted %d items in %v\n", insertCount, insertTime)

	// Check for existing and non-existing items.
	hits, falsePositives := 0, 0
	startCheck := time.Now()

	// Check for existing items.
	pipe = rdb.Pipeline()
	for i := 0; i < insertCount; i++ {
		pipe.Do(ctx, "BF.EXISTS", filterName, fmt.Sprintf("item%d", i))
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for existing items: %w", err)
	}
	for _, cmd := range cmds {
		if res, err := cmd.(*redis.Cmd).Int(); err == nil && res == 1 {
			hits++
		}
	}

	// Check for non-existing items (potential false positives).
	pipe = rdb.Pipeline()
	for i := 0; i < testSize; i++ {
		pipe.Do(ctx, "BF.EXISTS", filterName, fmt.Sprintf("fakeitem%d", i))
	}
	cmds, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for non-existing items: %w", err)
	}
	for _, cmd := range cmds {
		if res, err := cmd.(*redis.Cmd).Int(); err == nil && res == 1 {
			falsePositives++
		}
	}
	checkTime := time.Since(startCheck)

	fmt.Printf("  Hit rate: %d/%d, False positives: %d/%d (%.2f%%)\n",
		hits, insertCount, falsePositives, testSize, float64(falsePositives)/float64(testSize)*100)
	fmt.Printf("  Query time: %v\n", checkTime)
	fmt.Println("----------------------------------------------------")

	return nil
}


func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}

	configs := []testConfig{
		{1000, 0.01, 1000},
		{1000, 0.01, 10000},
		{100000, 0.001, 100000},
		{100000, 0.001, 120000},
		{1000000, 0.0001, 1000000},
	}

	for _, config := range configs {
		if err := testBloomFilter(ctx, rdb, config); err != nil {
			fmt.Printf("Error during test run for config %+v: %v\n", config, err)
		}
	}
}
