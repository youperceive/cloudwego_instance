#/bin/bash

kitexcall \
--idl-path ../../idl/user_account/user_account.thrift \
--method UserAccountService/Login \
--endpoint 127.0.0.1:8001 \
-d '{"target": "904413", "target_type": 1, "password": "123456"}'