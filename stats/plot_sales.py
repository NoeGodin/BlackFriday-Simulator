from plot_utils import load_csv_data, create_binned_plot

CSV_PATH = 'stats/sales_tracker.csv'
COLUMNS = ['simulation_id', 'map_name', 'tick', 'montant_vente', 'profit_total']
X_COLUMN = 'tick'
Y_COLUMN = 'profit_total'
GROUP_COLUMN = 'map_name'
INTERVAL = 50
TITLE = 'Evolution of mean profit by simulation ticks'
XLABEL = 'Simulation ticks'
YLABEL = 'Average total profit (€)'
OUTPUT_FILE = 'stats/profit_comparison.png'

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