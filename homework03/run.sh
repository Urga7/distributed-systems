#!/bin/bash
n=32
for ((i=n-1; i>=0; i--)); do
    ./app -pid=$i -n=$n -m=2 -k=3 &
done
wait

~/go/bin/GoVector --log_type shiviz --log_dir . --outfile Log-full.log

find . -name "Log-Process*.txt" -type f -delete
