#!/bin/bash

kitex \
    -module github.com/youperceive/cloudwego_instance/rpc/order \
    -service order \
    -I ../../idl \
    ../../idl/order/order.thrift