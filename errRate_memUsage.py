import pandas as pd
import matplotlib.pyplot as plt
import os

file_path = 'errRate_memUsage.csv'

if not os.path.exists(file_path):
    print(f"Error: The file '{file_path}' was not found.")
else:
    try:
        df = pd.read_csv(file_path)

        # Convert relevant columns to numeric types
        df['errorRate'] = pd.to_numeric(df['errorRate'])
        df['m'] = pd.to_numeric(df['m'])
        df['capacity'] = pd.to_numeric(df['capacity'])

        subset = df

        if subset.empty:
            print(f"No data found for capacity. Please check the CSV file or choose a different 'n'.")
        else:
            print(subset.head())
            plt.figure(figsize=(10, 6))
            plt.plot(subset['errorRate'], subset['m'] / (8 * 1024 * 1024))
            plt.gca().invert_xaxis()
            plt.xlabel('Expected Error Rate (p)')
            plt.ylabel('Total Memory Usage (KB)')
            plt.title(f'Bloom Filter Memory Usage vs. Expected Error Rate')
            plt.grid(True, which="both", ls="--", c='0.7')
            plt.tight_layout()

            # Save the plot
            plot_filename = f'errRate_memUsage.png'
            
            
            plt.savefig(plot_filename)
            print(f"Plot saved to {plot_filename}")

    except Exception as e:
        print(f"An error occurred while processing the CSV or plotting: {e}")
