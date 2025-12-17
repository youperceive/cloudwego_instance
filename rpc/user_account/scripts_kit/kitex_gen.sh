#!/bin/bash

kitex \
    -module github.com/youperceive/cloudwego_instance/rpc/user_account \
    -service user_account \
    -I ../../idl \
    ../../idl/UserAccountService/user.thrift