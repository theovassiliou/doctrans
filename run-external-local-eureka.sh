#!/bin/sh
docker container prune -f > /dev/null
docker run --name eureka -p 8761:8761 -d -it aista/eureka:latest
echo "Eureka server started at :8761"
