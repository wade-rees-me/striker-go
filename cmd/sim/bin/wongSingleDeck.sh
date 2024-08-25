#!/bin/bash

NUMBER_OF_ROUNDS=25000
NUMBER_OF_ROUNDS=$(cd $(dirname "$0"); RETURN_VALUE=$(source ./runSimulation.sh $1 $NUMBER_OF_ROUNDS); echo $RETURN_VALUE)

go run . -number-of-rounds $NUMBER_OF_ROUNDS -number-of-tables 1 -strategy-wong -table-single-deck -table-penetration 0.75 -table-blackjack-pays 3:2

