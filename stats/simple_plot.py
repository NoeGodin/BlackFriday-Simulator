import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

try:
    df = pd.read_csv('stats/sales_tracker.csv')
    if list(df.columns) != ['simulation_id', 'map_name', 'tick', 'montant_vente', 'profit_total']:
        print("No correct headers, using them")
        df = pd.read_csv('stats/sales_tracker.csv', header=None, names=['simulation_id', 'map_name', 'tick', 'montant_vente', 'profit_total'])
except:
    df = pd.DataFrame(columns=['simulation_id', 'map_name', 'tick', 'montant_vente', 'profit_total'])
plt.figure(figsize=(12, 8))

for map_name in df['map_name'].unique():
    map_data = df[df['map_name'] == map_name]
    
    # 100 ticks interval
    # TODO: maybe change interval ?
    max_tick = map_data['tick'].max()
    tick_bins = np.arange(0, max_tick + 100, 100)
    
    avg_profits = []
    tick_centers = []
    
    # mean for every tick interval
    for i in range(len(tick_bins) - 1):
        start_tick = tick_bins[i]
        end_tick = tick_bins[i + 1]
        center = (start_tick + end_tick) / 2
        
        # gather all simulation data
        interval_data = map_data[
            (map_data['tick'] >= start_tick) & 
            (map_data['tick'] < end_tick)
        ]
        
        if not interval_data.empty:
            avg_profit = interval_data['profit_total'].mean()
            avg_profits.append(avg_profit)
            tick_centers.append(center)
    
    # Trace line
    plt.plot(tick_centers, avg_profits, label=f'{map_name}', linewidth=2, marker='o', markersize=4)

plt.xlabel('Simulation ticks')
plt.ylabel('Average total profit (€)')
plt.title('Evolution of mean profit by simulation ticks')
plt.legend()
plt.grid(True, alpha=0.3)
plt.tight_layout()

# Save and display
plt.savefig('stats/profit_comparison.png', dpi=300, bbox_inches='tight')
plt.show()

print("Graph saved in stats/profit_comparison.png")