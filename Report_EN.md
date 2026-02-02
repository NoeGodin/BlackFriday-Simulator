# AI30 - Black Friday Simulation

## Project Description

The project we decided to develop is a supermarket simulation during a Black Friday or sales period. The main goal of this simulation is to reproduce customer behavior in a context of high traffic.

To represent a supermarket, we adopted a top-down view. On our interface, we represented the following elements:
- Shelves: containing products that agents will go to fetch.
- Cash registers: for our supermarket to record products purchased by our agents.
- Doors: where agents appear at regular intervals and where they go once they have paid for their products.

Agents are customers who will collect up to x items simultaneously and have several behaviors:
- Collaborative: these types of agents will go fetch the products they need.
- Selfish: act like collaborative agents, if they look for a product that a nearby agent possesses, they can steal their products.
- Resentful: if they get robbed, they will try to steal someone else to recover their product.
- Security Guard: these types of agents have the particularity of moving randomly in the supermarket and help reduce inter-agent theft around them.

In order to better understand the actions taking place in our simulation, by clicking on elements on the map (agents, shelves, ...), information related to it will be displayed in a HUD (Heads up display). Furthermore, when an agent picks up an item, they are highlighted in green.

## Question to Answer

The question we seek to answer with our simulation is: what is the best store layout to make the most profit, as quickly as possible?

## Response Indicators

The indicators used to answer this question are as follows:
- Supermarket profit on a layout: how much the supermarket can profit in its store from a given stock.
- Time agents take to make their purchases: the best layout involves the time agents take to quickly find their products and pay for them quickly before leaving.
- Number of inter-agent collisions: to consider the well-being of agents, the number of collisions is a number that negatively affects the store layout (aisles too tight, difficult access to shelves or cash registers).

## Architecture Employed

Regarding the architecture we employ, we use the Ebitengine graphics library for our window.

### Internal Packages Architecture

**pkg/constants**: centralizes all constants of our application (cell size, social force coefficient, agent vision size, ...).

**pkg/graphics**: handles the user interface, including interactions with the interface, agent animations, image management, ...

**pkg/hud**: manages the HUD display on the window.

**pkg/logger**: allows displaying information to debug our application.

**pkg/map**: handles the simulation map.

**pkg/pathfinding**: manages agent pathfinding using the A* algorithm.

**pkg/shopping**: generates and reads from files the shopping lists for our agents.

**pkg/simulation**: contains the core of our simulation, including agent management, concurrent access, collisions, visions...

**pkg/utils**: utility folder to manage shared types and mathematical calculations.

### Commands

**cmd/blackfriday**: launches our simulation by reading a file for a map to generate, a JSON list for stocks, and a JSON list for the agent shopping list.
**cmd/generate_shopping_lists_from_stocks**: generates a shopping list for all agents from stock data. Can make, if desired, the shopping list deterministic across different simulations.
**cmd/map_generator**: randomly generates maps including a given size, number of doors, cash registers, shelves, walls.

## Parameters

The parameters we propose to change are the following. These are to be filled in a `.env` file. A `.env.example` file is in the repository with said parameters:
- NB_AGENTS (the number of agents)
- AGENT_SPEED (the agent movement speed)
- AGENT_MAX_SHOPPING_LIST (the maximum number of products an agent may want to fetch)

If these parameters are not filled, the simulation assigns arbitrary values.

## How to Run the Application

To run the application, execute the following commands:

#### Installing Dependencies
```
sudo apt install -y libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev
sudo apt install -y libgl1-mesa-dev xorg-dev
```

#### Configuration
```
cp .env.example .env
# Modify .env according to your needs (number of agents, speed, etc.)

go run cmd/generate_shopping_lists_from_stocks/main.go
# To generate shopping lists for agents, in order to compare multiple simulations with the same data.
```

#### Execution
```
go mod download
go run cmd/blackfriday/main.go
```

#### Sales Analysis
```
./stats/plot.sh
# Runs a script that will display the chart from data extracted from simulations.
```
