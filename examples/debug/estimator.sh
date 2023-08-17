#!/bin/bash

# Initialize total time to zero
total_time=0.0

# Number of runs
runs=20

# Difficulty value
difficulty=6

# Run the command for the specified number of times
for i in $(seq 1 $runs); do
    # Capture the real time taken by the command using the shell's built-in 'time' command
    # We'll then filter and extract only the real elapsed time value using awk
    time_output=$( { time go run main.go $difficulty; } 2>&1 | grep real | awk '{print $2}' | sed 's/m/ * 60 + /' | sed 's/s//' | bc -l)
    total_time=$(printf "%.4f + %.4f\n" "$total_time" "$time_output" | bc -l)
done

# Calculate average time
average_time=$(printf "%.4f / %d\n" "$total_time" "$runs" | bc -l)

echo "Average run time for difficulty $difficulty over $runs runs: $average_time seconds"
