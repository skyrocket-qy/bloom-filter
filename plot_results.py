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
        df['falsePositiveRate'] = pd.to_numeric(df['falsePositiveRate'])
        df['capacity'] = pd.to_numeric(df['capacity'])

        plt.figure(figsize=(12, 7))

        # Plot falsePositiveRate vs errorRate for each capacity
        for capacity in df['capacity'].unique():
            subset = df[df['capacity'] == capacity]
            plt.plot(subset['errorRate'], subset['falsePositiveRate'], marker='o', linestyle='-', label=f'Capacity: {capacity}')

        plt.xscale('log') # Error rates are often logarithmic
        plt.xlabel('Expected Error Rate (p)')
        plt.ylabel('Actual False Positive Rate (%)')
        plt.title('Bloom Filter False Positive Rate vs. Expected Error Rate')
        plt.grid(True, which="both", ls="--", c='0.7')
        plt.legend()
        plt.tight_layout()

        # Save the plot
        plot_filename = 'bloom_filter_plot.png'
        plt.savefig(plot_filename)
        print(f"Plot saved to {plot_filename}")

        plt.show()

    except Exception as e:
        print(f"An error occurred while processing the CSV or plotting: {e}")
