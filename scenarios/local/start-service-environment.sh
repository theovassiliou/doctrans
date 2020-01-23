#!/bin/sh
echo `pwd`
rm -f .env
echo "HOST=`hostname`" >> .env
echo "IP=`/sbin/ifconfig en0 | grep "inet " | awk -F' ' '{print $2}' | awk '{print $1}'`" >> .env
docker-compose -f scenarios/local/docker-compose.yaml up -d 
