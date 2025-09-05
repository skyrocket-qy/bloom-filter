# Bloom Filter Performance Analysis

This project analyzes the performance characteristics of a Bloom filter implementation that uses Redis's `rebloom` module. It provides tools to measure and visualize the trade-offs between memory usage, error rate, and performance.

## What is a Bloom Filter?

A Bloom filter is a space-efficient probabilistic data structure that is used to test whether an element is a member of a set. It is "probabilistic" because it can return two types of answers:

*   **"Definitely not in the set."**
*   **"Possibly in the set."**

This means that false positive matches are possible, but false negatives are not. In other words, the filter can tell you that an item *might* be in the set when it isn't, but it will never tell you an item *isn't* in the set when it is.

The main advantages of Bloom filters are their speed and memory efficiency. They are ideal for applications where a small rate of false positives is acceptable, such as cache filtering, network routing, and spell checking.

For a deeper dive into other types of filters and when to use them, see [this guide](./filter.md).

## Prerequisites

To run this project, you will need:

*   [Go](https://golang.org/)
*   [Docker](https://www.docker.com/)
*   [Python 3](https://www.python.org/)
*   The following Python packages: `pandas`, `matplotlib`

You can install the required Python packages using pip:
```bash
pip install pandas matplotlib
```

## How to Run

1.  **Start the RedisBloom server:**

    This project requires a Redis server with the RedisBloom module. The included `Makefile` provides a convenient way to start a pre-configured Docker container.

    ```bash
    make start-redis
    ```

2.  **Run the analysis:**

    The main program is written in Go. It will run a series of tests against the Redis server and generate CSV files with the results in the `out/` directory.

    You can run all tests at once:
    ```bash
    make run
    ```
    This is equivalent to `go run main.go -test=all`.

    Alternatively, you can run specific tests using the `-test` flag:
    *   `mem`: Generates data for error rate vs. memory usage.
    *   `fp`: Generates data for the actual false positive rate vs. the number of inserted items.
    *   `time`: Generates data for error rate vs. check time.

    Example:
    ```bash
    go run main.go -test=mem
    ```

3.  **Generate the plots:**

    After generating the CSV data, you can use the Python scripts in the `scripts/` directory to create plots.

    ```bash
    python3 scripts/errRate_memUsage.py
    python3 scripts/realAmount_fpRate.py
    python3 scripts/errRate_checkTime.py
    ```

    The generated plots will be saved as PNG images in the `out/` directory.

4.  **Stop the RedisBloom server:**

    When you are finished, you can stop and remove the Redis container:

    ```bash
    make stop-redis
    ```

5.  **Clean up generated files:**

    To remove the `out/` directory and all its contents:
    ```bash
    make clean
    ```

## Understanding the Results

The analysis generates three different plots:

*   **Error Rate vs. Memory Usage (`errRate_memUsage.png`):** This plot shows the relationship between the desired false positive rate (`p`) and the amount of memory (`m`) required. You will see that requiring a lower error rate (fewer false positives) demands significantly more memory.

*   **Real Amount vs. False Positive Rate (`realAmount_fpRate.png`):** This plot demonstrates how the false positive rate is affected by the number of items actually inserted into the filter, especially when it exceeds the filter's original capacity (`n`).

*   **Error Rate vs. Check Time (`errRate_checkTime.png`):** This plot shows how the time it takes to check for an item's existence is affected by the filter's configuration.

## Key Formulas

The behavior of a Bloom filter is governed by these two formulas, which calculate the optimal number of bits (`m`) and hash functions (`k`) for a given number of items (`n`) and a desired false positive rate (`p`).

$$
m = -\\frac{n \\ln p}{(\\ln 2)^2}
$$

$$
k = \\frac{m}{n} \\ln 2
$$

Where:
*   `m` = total number of bits in the filter's bit array
*   `n` = expected number of distinct items
*   `p` = target false positive rate
*   `k` = number of hash functions
