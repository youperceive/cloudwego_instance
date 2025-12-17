#!/bin/bash

kitex \
    -module github.com/youperceive/cloudwego_instance/rpc/verify_code \
    -service verify_code \
    -I ../../idl \
    ../../idl/verify_code/verify_code.thrift