#!/bin/bash

cd "$(dirname "$0")/.."

# Create venv is not exist
if [ ! -d "venv" ]; then
    echo "Creating venv..."
    python3 -m venv venv
fi

source venv/bin/activate
echo "Verifying dependencies..."
pip install matplotlib pandas > /dev/null 2>&1

echo "Generating graphic..."
python stats/simple_plot.py