#!/bin/sh
rm -f .env
echo "HOST=`hostname`" >> .env
echo "IP=`/sbin/ifconfig en0 | grep "inet " | awk -F' ' '{print $2}' | awk '{print $1}'`" >> .env
docker-compose -f docker-compose.yaml up -d 
