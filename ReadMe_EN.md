# Black Friday Simulator

## Installation

**Linux (Ubuntu/Debian)**: Install system dependencies:
```bash
sudo apt install -y libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev
sudo apt install -y libgl1-mesa-dev xorg-dev
```

**Configuration:**
```bash
cp .env.example .env
# Modify .env according to your needs (number of agents, speed, etc.)
```

**Execution:**
```bash
go mod download
go run cmd/blackfriday/main.go
```

**Sales Analysis:** Simulations automatically record every sale with a timestamp in `stats/sales_tracker.csv`. To generate a chart comparing the performance of different maps, use:
```bash
./stats/plot.sh
```

**Deterministic Shopping Lists:** For fair comparisons between different store configurations, generate predefined shopping lists based on stocks:
```bash
go run cmd/generate_shopping_lists_from_stocks/main.go
```
This creates a `maps/store/shopping_lists.json` file with deterministic lists for each agent, guaranteeing that all agents will have the same lists in every simulation.

## Problem Statement

Store's point of view: How to maximize sales?

## Functioning

On the graphical interface, one can place:
- shelves (obstacles)
- items (goal)
- an entrance / exit

Agents would be customers who could collect up to x items simultaneously and have several behaviors:
- collaborative: leave the item to an agent within a distance range without an item around them unless they are also collaborative
- competitive: will take the optimal path to enter and exit with 1 item
- selfish: will steal up to 3 items from other agents if no other agents with items will go look for one directly
- resentful: competitive, if they get robbed, will steal from someone else other than their thief

We could also have agents that restock shelves.
See how we manage item quantities.
