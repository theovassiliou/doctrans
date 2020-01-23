#!/bin/bash
for file in scenarios/global/qds_*.sh; do 
    if [ -f "$file" ]; then 
        killall -9 `basename -s .sh $file` 
    fi 
done