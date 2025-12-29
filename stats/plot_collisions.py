from plot_utils import load_csv_data, create_binned_plot

CSV_PATH = 'stats/collision_tracker.csv'
COLUMNS = ['simulation_id', 'map_name', 'tick', 'collision_count', 'total_collisions']
X_COLUMN = 'tick'
Y_COLUMN = 'total_collisions'
GROUP_COLUMN = 'map_name'
INTERVAL = 50
TITLE = 'Evolution of mean collisions by simulation ticks'
XLABEL = 'Simulation ticks'
YLABEL = 'Average total collisions'
OUTPUT_FILE = 'stats/collision_comparison.png'

collision_df = load_csv_data(CSV_PATH, COLUMNS)
create_binned_plot(
    df=collision_df,
    x_col=X_COLUMN,
    y_col=Y_COLUMN,
    group_col=GROUP_COLUMN,
    interval=INTERVAL,
    title=TITLE,
    xlabel=XLABEL,
    ylabel=YLABEL,
    output_file=OUTPUT_FILE
)