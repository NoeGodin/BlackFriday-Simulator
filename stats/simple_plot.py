import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

try:
    df = pd.read_csv('stats/sales_tracker.csv')
    if list(df.columns) != ['simulation_id', 'map_name', 'temps_relatif_sec', 'timestamp', 'montant_vente', 'profit_total']:
        print("No correct headers, using them")
        df = pd.read_csv('stats/sales_tracker.csv', header=None, names=['simulation_id', 'map_name', 'temps_relatif_sec', 'timestamp', 'montant_vente', 'profit_total'])
except:
    df = pd.DataFrame(columns=['simulation_id', 'map_name', 'temps_relatif_sec', 'timestamp', 'montant_vente', 'profit_total'])
plt.figure(figsize=(12, 8))

for map_name in df['map_name'].unique():
    map_data = df[df['map_name'] == map_name]
    
    # 10 seconds time interval
    # TODO: maybe change interval ?
    max_time = map_data['temps_relatif_sec'].max()
    time_bins = np.arange(0, max_time + 10, 10)
    
    avg_profits = []
    time_centers = []
    
    # mean for every time interval
    for i in range(len(time_bins) - 1):
        start_time = time_bins[i]
        end_time = time_bins[i + 1]
        center = (start_time + end_time) / 2
        
        # gather all simulation data
        interval_data = map_data[
            (map_data['temps_relatif_sec'] >= start_time) & 
            (map_data['temps_relatif_sec'] < end_time)
        ]
        
        if not interval_data.empty:
            avg_profit = interval_data['profit_total'].mean()
            avg_profits.append(avg_profit)
            time_centers.append(center)
    
    # Trace line
    plt.plot(time_centers, avg_profits, label=f'{map_name}', linewidth=2, marker='o', markersize=4)

plt.xlabel('Time (seconds)')
plt.ylabel('Average total profit (€)')
plt.title('Evolution of mean profit by time')
plt.legend()
plt.grid(True, alpha=0.3)
plt.tight_layout()

# Save and display
plt.savefig('stats/profit_comparison.png', dpi=300, bbox_inches='tight')
plt.show()

print("Graph saved in stats/profit_comparison.png")