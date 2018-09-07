#!/bin/bash
echo "Running multiple instances."

for i in {1..10}
do
  DEBUG=1 fn call tenDynamos /droplet$i > call_$i.log &
  echo "finished round $i"
done

COLOUR="yellow"
echo "The flower is $COLOUR"
