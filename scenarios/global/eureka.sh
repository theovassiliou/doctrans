#!/bin/sh
# TODO Check for running container
docker container prune -f > /dev/null
docker run --name eureka -p 8761:8761 -d -it aista/eureka:latest
# TODO Add a test (curl) that server is up and accessible
echo "Eureka server started at :8761"
