from plot_utils import load_csv_data, create_binned_plot

CSV_PATH = 'stats/aggressiveness_tracker.csv'
COLUMNS = ["simulation_id", "map_name", "tick", "aggressiveness"]
X_COLUMN = 'tick'
Y_COLUMN = 'aggressiveness'
GROUP_COLUMN = 'map_name'
INTERVAL = 50
TITLE = 'Evolution of mean aggressiveness by simulation ticks'
XLABEL = 'Simulation ticks'
YLABEL = 'mean aggressiveness'
OUTPUT_FILE = 'stats/aggressiveness_comparison.png'

sales_df = load_csv_data(CSV_PATH, COLUMNS)
create_binned_plot(
    df=sales_df,
    x_col=X_COLUMN,
    y_col=Y_COLUMN,
    group_col=GROUP_COLUMN,
    interval=INTERVAL,
    title=TITLE,
    xlabel=XLABEL,
    ylabel=YLABEL,
    output_file=OUTPUT_FILE
)