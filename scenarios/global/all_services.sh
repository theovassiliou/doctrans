#!/bin/bash
for file in scenarios/global/qds_*.sh; do 
    if [ -f "$file" ]; then 
        $file 
    fi 
done