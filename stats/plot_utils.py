import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

def load_csv_data(csv_path, columns):
    try:
        df = pd.read_csv(csv_path)
        if list(df.columns) != columns:
            print(f"No correct headers for {csv_path}, using manual names")
            df = pd.read_csv(csv_path, header=None, names=columns)
        return df
    except:
        print(f"Could not load {csv_path}, creating empty DataFrame")
        return pd.DataFrame(columns=columns)

def create_binned_plot(df, x_col, y_col, group_col, interval, title, xlabel, ylabel, output_file):
    plt.figure(figsize=(12, 8))
    
    if not df.empty:
        for group_name in df[group_col].unique():
            group_data = df[df[group_col] == group_name]
            print(f"Processing {group_name}, data points: {len(group_data)}")
            
            if len(group_data) > 0:
                # Create bins for averaging
                max_x = group_data[x_col].max()
                x_bins = np.arange(0, max_x + interval, interval)
                
                avg_values = []
                x_centers = []
                
                # Calculate average for each interval
                for i in range(len(x_bins) - 1):
                    start_x = x_bins[i]
                    end_x = x_bins[i + 1]
                    center = (start_x + end_x) / 2
                    
                    # Get data in this interval
                    interval_data = group_data[
                        (group_data[x_col] >= start_x) & 
                        (group_data[x_col] < end_x)
                    ]
                    
                    if not interval_data.empty:
                        avg_value = interval_data[y_col].mean()
                        avg_values.append(avg_value)
                        x_centers.append(center)
                
                # Plot the line
                if len(x_centers) > 0:
                    plt.plot(x_centers, avg_values, label=f'{group_name}', 
                            linewidth=2, marker='o', markersize=4)
    else:
        print("DataFrame is empty")
    
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.title(title)
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()

    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    print(f"Graph saved in {output_file}")
    plt.close()