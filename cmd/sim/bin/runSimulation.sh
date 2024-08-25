#!/bin/bash

NUMBER_OF_ROUNDS=$2

# Check if a command line argument is provided
if [ -n "$1" ]; then
    NUMBER_OF_ROUNDS="$1"
fi

echo $NUMBER_OF_ROUNDS
