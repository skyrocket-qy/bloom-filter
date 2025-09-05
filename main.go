package main

import (
	"context"
	"encoding/csv" // New import
	"fmt"
	"log"
	"math" // New import
	"os"
	"os/exec"
	"strconv" // New import for converting numbers to string
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

// testConfig holds the configuration for a single Bloom filter test run.
type testConfig struct {
	capacity  int
	errorRate float64
}

// testBloomFilter runs a single Bloom filter test with the given configuration.
// It now accepts a csv.Writer to write results.
func testBloomFilter(ctx context.Context, rdb *redis.Client, config testConfig, actualInsertCount int, writer *csv.Writer) error {
	const testSize = 1000

	// Calculate Bloom filter parameters.
	m, k := BloomFilterParams(config.capacity, config.errorRate)
	// fmt.Printf("Testing n=%d, p=%.4f -> m=%d, k=%d\n", config.capacity, config.errorRate, m, k) // Remove this line

	filterName := fmt.Sprintf("filter_n%d_p%.4f", config.capacity, config.errorRate)
	defer rdb.Del(ctx, filterName) // Ensure cleanup

	rdb.Del(ctx, filterName) // Ensure the filter doesn't exist

	// Reserve the filter.
	if err := rdb.Do(ctx, "BF.RESERVE", filterName, config.errorRate, config.capacity).Err(); err != nil {
		return fmt.Errorf("failed to reserve bloom filter: %w", err)
	}

	// Insert items using a pipeline for efficiency.
	pipe := rdb.Pipeline()
	insertCount := actualInsertCount
	startInsert := time.Now()
	for i := range insertCount {
		pipe.Do(ctx, "BF.ADD", filterName, fmt.Sprintf("item%d", i))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert items: %w", err)
	}
	insertTime := time.Since(startInsert)

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
		if res, err := cmd.(*redis.Cmd).Result(); err == nil {
			if val, ok := res.(bool); ok && val {
				hits++
			}
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
		if res, err := cmd.(*redis.Cmd).Result(); err == nil {
			if val, ok := res.(bool); ok && val {
				falsePositives++
			}
		}
	}
	checkTime := time.Since(startCheck)

	falsePositiveRate := float64(falsePositives) / float64(testSize)
	// Write results to CSV
	record := []string{
		strconv.Itoa(config.capacity),
		fmt.Sprintf("%.9f", config.errorRate),
		fmt.Sprintf("%.4f", float64(m)/1024),
		strconv.Itoa(k),
		insertTime.String(),
		checkTime.String(),
		strconv.Itoa(hits),
		strconv.Itoa(insertCount),
		strconv.Itoa(falsePositives),
		strconv.Itoa(testSize),
		fmt.Sprintf("%.4f", falsePositiveRate),
	}
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write CSV record: %w", err)
	}
	writer.Flush() // Ensure data is written immediately

	return nil
}

func TestRealAmount_fpRate(ctx context.Context, rdb *redis.Client) {
	// Open CSV file for writing
	file, err := os.OpenFile("realAmount_fpRate.csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Failed to open CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all buffered writes are flushed before closing

	// Check if file is empty to write header
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		header := []string{
			"capacity", "errorRate", "m", "k", "insertTime", "checkTime",
			"hits", "insertCount", "falsePositives", "testSize", "falsePositiveRate",
		}
		if err := writer.Write(header); err != nil {
			fmt.Printf("Failed to write CSV header: %v\n", err)
			return
		}
	}

	// Test for varying realAmount (insertCount) for a fixed n and p
	fixedN := 10000
	fixedP := 0.001

	// Define a range of realAmount values to test
	realAmounts := []int{
		1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000,
		11000, 12000, 13000, 14000, 15000, 16000, 17000, 18000, 19000, 20000,
	} // Example values

	for _, realAmount := range realAmounts {
		config := testConfig{capacity: fixedN, errorRate: fixedP}
		if err := testBloomFilter(ctx, rdb, config, realAmount, writer); err != nil {
			fmt.Printf("Error during test run for varying realAmount (n=%d, p=%.4f, realAmount=%d): %v\n", fixedN, fixedP, realAmount, err)
		}
	}

	cmd := exec.Command("python3", "plot_realAmount_fpRate.py")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	log.Println("Python script executed successfully.")
}

func TestErrRateVsMemUsage(ctx context.Context, rdb *redis.Client) {
	file, err := os.OpenFile("errRate_memUsage.csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Failed to open CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all buffered writes are flushed before closing

	// Check if file is empty to write header
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		header := []string{
			"capacity", "errorRate", "m", "k", "insertTime", "checkTime",
			"hits", "insertCount", "falsePositives", "testSize", "falsePositiveRate",
		}
		if err := writer.Write(header); err != nil {
			fmt.Printf("Failed to write CSV header: %v\n", err)
			return
		}
	}

	ns := []int{10e4}
	testConfigs := []testConfig{}
	for _, n := range ns {
		for p := 0.1; p >= 0.00001; {
			testConfigs = append(testConfigs, testConfig{n, p})
			p /= 1.05
		}
	}

	for _, config := range testConfigs {
		// Pass the writer to the testBloomFilter function
		if err := testBloomFilter(ctx, rdb, config, config.capacity, writer); err != nil {
			fmt.Printf("Error during test run for config %+v: %v\n", config, err)
		}
	}

	cmd := exec.Command("python3", "plot_errRate_memUsage.py")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	log.Println("Python script executed successfully.")
}

func TestErrRateVsCheckTime(ctx context.Context, rdb *redis.Client) {
	file, err := os.OpenFile("errRate_checkTime.csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Failed to open CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all buffered writes are flushed before closing

	// Check if file is empty to write header
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		header := []string{
			"capacity", "errorRate", "m", "k", "insertTime", "checkTime",
			"hits", "insertCount", "falsePositives", "testSize", "falsePositiveRate",
		}
		if err := writer.Write(header); err != nil {
			fmt.Printf("Failed to write CSV header: %v\n", err)
			return
		}
	}

	ns := []int{10e4}
	testConfigs := []testConfig{}
	for _, n := range ns {
		for p := 0.1; p >= 0.00000001; {
			testConfigs = append(testConfigs, testConfig{n, p})
			p /= 10
		}
	}

	for _, config := range testConfigs {
		// Pass the writer to the testBloomFilter function
		if err := testBloomFilter(ctx, rdb, config, config.capacity, writer); err != nil {
			fmt.Printf("Error during test run for config %+v: %v\n", config, err)
		}
	}

	cmd := exec.Command("python3", "plot_errRate_checkTime.py")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	log.Println("Python script executed successfully.")
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

	TestErrRateVsMemUsage(ctx, rdb)
	TestRealAmount_fpRate(ctx, rdb)
	TestErrRateVsCheckTime(ctx, rdb)
}
