#/bin/bash

kitexcall \
--idl-path ../../idl/user_account/user_account.thrift \
--method UserAccountService/Register \
--endpoint 127.0.0.1:8001 \
-d '{"username": "sxkane", "target": "sxshenxu", "register_type": 1, "password": "123456", "captcha": "877495"}'