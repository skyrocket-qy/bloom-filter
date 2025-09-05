import pandas as pd
import matplotlib.pyplot as plt
import os

file_path = 'realAmount_fpRate.csv'

if not os.path.exists(file_path):
    print(f"Error: The file '{file_path}' was not found.")
else:
    try:
        df = pd.read_csv(file_path)

        # Convert relevant columns to numeric types
        df['capacity'] = pd.to_numeric(df['capacity'])
        df['errorRate'] = pd.to_numeric(df['errorRate'])
        df['insertCount'] = pd.to_numeric(df['insertCount'])
        df['falsePositiveRate'] = pd.to_numeric(df['falsePositiveRate'])

        subset = df

        if subset.empty:
            print(f"No data found. Please check the CSV file or choose different values.")
        else:
            print(subset.head())
            plt.figure(figsize=(10, 6))
            plt.plot(subset['insertCount'], subset['falsePositiveRate'], marker='o', linestyle='-')

            plt.xlabel('Real Amount (Number of Inserted Items)')
            plt.ylabel('False Positive Rate (%)')
            plt.title(f'Bloom Filter False Positive Rate vs. Real Amount')
            plt.grid(True, which="both", ls="--", c='0.7')
            plt.tight_layout()

            plot_filename = f'realAmount_fpRate.png'
            plt.savefig(plot_filename)
            print(f"Plot saved to {plot_filename}")

    except Exception as e:
        print(f"An error occurred while processing the CSV or plotting: {e}")
