#!/bin/bash

docker exec -it order-mongo mongosh -u root -p root123456 --authenticationDatabase admin