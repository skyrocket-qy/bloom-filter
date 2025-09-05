import pandas as pd
import matplotlib.pyplot as plt
import os

file_path = 'bloom_filter_results.csv'

if not os.path.exists(file_path):
    print(f"Error: The file '{file_path}' was not found.")
else:
    try:
        df = pd.read_csv(file_path)

        # Convert relevant columns to numeric types
        df['errorRate'] = pd.to_numeric(df['errorRate'])
        df['m'] = pd.to_numeric(df['m'])
        df['capacity'] = pd.to_numeric(df['capacity'])

        # Choose a fixed capacity (n)
        # Let's pick the first unique capacity found in the data, or a specific one like 10000
        fixed_n = df['capacity'].unique()[0] # Or set to 10000, 100000 etc.
        print(f"Plotting memory usage for fixed capacity (n) = {fixed_n}")

        subset = df[df['capacity'] == fixed_n]

        if subset.empty:
            print(f"No data found for capacity (n) = {fixed_n}. Please check the CSV file or choose a different 'n'.")
        else:
            print(subset.head())
            plt.figure(figsize=(10, 6))
            plt.plot(subset['errorRate'], subset['m'])

            plt.xlabel('Expected Error Rate (p)')
            plt.ylabel('Total Memory Usage (MB)')
            plt.title(f'Bloom Filter Memory Usage vs. Expected Error Rate (n={fixed_n})')
            plt.grid(True, which="both", ls="--", c='0.7')
            plt.tight_layout()

            # Save the plot
            plot_filename = f'memory_usage_n{fixed_n}.png'
            # plt.yticks(range(int(subset['m'].min()), int(subset['m'].max()) + 1, 100000))
            plt.gca().invert_xaxis()
            plt.savefig(plot_filename)
            print(f"Plot saved to {plot_filename}")

    except Exception as e:
        print(f"An error occurred while processing the CSV or plotting: {e}")
