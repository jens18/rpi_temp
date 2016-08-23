#!/bin/bash

set -x

# start Docker container instance in daemon mode with the name 'rpi_temp'
docker run -p 9090:9090 --name rpi_temp -d rpi_temp

# run test
bash ./curl_loop.sh

# display Docker container log file
docker logs rpi_temp

# shutdown Docker
docker stop rpi_temp

# remove Docker instance
docker rm rpi_temp



