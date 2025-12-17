#!/bin/bash

kitex \
    -module github.com/cloudwego_instance/rpc/verify_code \
    -service verify_code \
    -I ../../idl \
    ../../idl/VerifyCodeService/captcha.thrift