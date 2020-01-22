#!/bin/sh
rm -f .env
echo "HOST=`hostname`" >> .env
docker-compose -f docker-compose.yaml up -d 
