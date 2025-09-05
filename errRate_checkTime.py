import pandas as pd
import matplotlib.pyplot as plt
import os
import numpy as np
import numpy as np

file_path = 'out/errRate_checkTime.csv'

if not os.path.exists(file_path):
    print(f"Error: The file '{file_path}' was not found.")
else:
    try:
        df = pd.read_csv(file_path)

        # Convert relevant columns to numeric types
        df['errorRate'] = pd.to_numeric(df['errorRate'])
        df['capacity'] = pd.to_numeric(df['capacity'])

        # Convert checkTime from string (e.g., "123.456ms") to milliseconds (float)
        def parse_time_to_ms(time_str):
            if isinstance(time_str, str) and time_str.endswith('ms'):
                try:
                    return float(time_str[:-2])
                except ValueError:
                    return None
            return None # Or handle other units if necessary

        df['checkTimeMs'] = df['checkTime'].apply(parse_time_to_ms)

        subset = df

        if subset.empty:
            print(f"No data found. Please check the CSV file or choose a different 'n'.")
        else:
            plt.figure(figsize=(10, 6))
            plt.plot(subset['errorRate'], subset['checkTimeMs'], marker='o', linestyle='-')
            plt.gca().invert_xaxis()
            plt.xscale('log')
            plt.xlabel('Expected Error Rate (p)')
            plt.ylabel('Check Time (ms)')
            plt.title(f'Bloom Filter Check Time vs. Expected Error Rate')
            plt.grid(True, which="both", ls="--", c='0.7')
            plt.tight_layout()

            plot_filename = f'out/errRate_checkTime.png'
            plt.savefig(plot_filename)
            print(f"Plot saved to {plot_filename}")

    except Exception as e:
        print(f"An error occurred while processing the CSV or plotting: {e}")
