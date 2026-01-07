#!/bin/bash

cd "$(dirname "$0")/.."

# Pick a Python executable available on the system
PYTHON_CMD=""
if command -v python3 >/dev/null 2>&1; then
    PYTHON_CMD="python3"
elif command -v python >/dev/null 2>&1; then
    PYTHON_CMD="python"
elif command -v py >/dev/null 2>&1; then
    PYTHON_CMD="py -3"
else
    echo "Error: Python not found (python3/python/py)." >&2
    exit 1
fi

# Create venv is not exist
if [ ! -d "venv" ]; then
    echo "Creating venv..."
    $PYTHON_CMD -m venv venv
fi

# Use the venv Python directly (works on Windows + Linux/macOS)
VENV_PY=""
if [ -f "venv/Scripts/python.exe" ]; then
    VENV_PY="venv/Scripts/python.exe"
elif [ -f "venv/bin/python" ]; then
    VENV_PY="venv/bin/python"
else
    echo "Error: venv python not found. Try deleting ./venv and re-running." >&2
    exit 1
fi

echo "Verifying dependencies..."
$VENV_PY -m pip install --disable-pip-version-check matplotlib pandas

echo "Generating sales graphic..."
$VENV_PY stats/plot_sales.py

echo "Generating collisions graphic..."
$VENV_PY stats/plot_collisions.py

echo "Generating aggressiveness graphic..."
$VENV_PY stats/plot_aggressiveness.py

echo "Generating steals graphic..."
$VENV_PY stats/plot_steals.py
