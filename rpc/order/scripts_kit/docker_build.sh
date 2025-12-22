#!/bin/bash

docker build -f Dockerfile -t user-account-service .

# docker stop $(docker ps -aq)
# docker rm $(docker ps -aq)