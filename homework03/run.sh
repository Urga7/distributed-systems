#!/bin/bash
for i in {4..0}; do
    ./app -pid=$i -n=5 -m=3 -k=3 > "output_$i.log" 2>&1 &
done
wait 