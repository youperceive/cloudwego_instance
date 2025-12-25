#!/bin/bash

docker build -f Dockerfile -t order-service .

# docker stop $(docker ps -aq)
# docker rm $(docker ps -aq)