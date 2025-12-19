#/bin/bash

kitexcall \
--idl-path ../../idl/user_account/user_account.thrift \
--method UserAccountService/Register \
--endpoint 127.0.0.1:8001 \
-d '{"username": "sxkane", "target": "sxshenxu", "target_type": 1, "password": "123456", "captcha": "645747"}'