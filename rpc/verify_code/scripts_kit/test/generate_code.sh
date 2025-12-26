#!/bin/bash

kitexcall \
--idl-path ../../idl/verify_code/verify_code.thrift \
--method VerifyCodeService/GenerateCaptcha \
--endpoint 127.0.0.1:8000 \
-d '{"type": 1, "proj": "user-account-service", "biz_type": "login", "target": "904413"}'