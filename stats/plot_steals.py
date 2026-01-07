from plot_utils import load_csv_data
import matplotlib.pyplot as plt
import numpy as np

CSV_PATH = 'stats/steal_tracker.csv'
COLUMNS = ["simulation_id", "map_name", "tick", "stealer_id", "victim_id", "item_name"]
X_COLUMN = 'tick'
GROUP_COLUMN = 'map_name'
INTERVAL = 50
TITLE = 'Cumulative steals over simulation ticks'
XLABEL = 'Simulation ticks'
YLABEL = 'Total steals'
OUTPUT_FILE = 'stats/steal_comparison.png'


def create_cumulative_binned_count_plot(df, x_col, group_col, interval, title, xlabel, ylabel, output_file):
    plt.figure(figsize=(12, 8))

    if df.empty:
        print("DataFrame is empty")
    else:
        for group_name in df[group_col].unique():
            group_data = df[df[group_col] == group_name]
            print(f"Processing {group_name}, data points: {len(group_data)}")

            if len(group_data) == 0:
                continue

            group_data = group_data.sort_values(by=x_col)
            max_x = group_data[x_col].max()
            x_bins = np.arange(0, max_x + interval, interval)

            cumulative_counts = []
            x_centers = []
            running_total = 0

            for i in range(len(x_bins) - 1):
                start_x = x_bins[i]
                end_x = x_bins[i + 1]
                center = (start_x + end_x) / 2

                interval_count = len(group_data[
                    (group_data[x_col] >= start_x) &
                    (group_data[x_col] < end_x)
                ])

                running_total += interval_count
                cumulative_counts.append(running_total)
                x_centers.append(center)

            if len(x_centers) > 0:
                plt.plot(x_centers, cumulative_counts, label=f'{group_name}', linewidth=2)

    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.title(title)
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()

    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    print(f"Graph saved in {output_file}")
    plt.close()


steal_df = load_csv_data(CSV_PATH, COLUMNS)

# Ensure tick is numeric when the CSV is loaded without headers
try:
    steal_df[X_COLUMN] = steal_df[X_COLUMN].astype(int)
except Exception:
    pass

create_cumulative_binned_count_plot(
    df=steal_df,
    x_col=X_COLUMN,
    group_col=GROUP_COLUMN,
    interval=INTERVAL,
    title=TITLE,
    xlabel=XLABEL,
    ylabel=YLABEL,
    output_file=OUTPUT_FILE
)
