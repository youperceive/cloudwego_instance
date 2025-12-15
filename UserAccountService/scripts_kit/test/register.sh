#/bin/bash

kitexcall \
--idl-path idl/user.thrift \
--method UserAccountService/Register \
--endpoint 127.0.0.1:8001 \
-d '{"username": "AlfredGit", "target": "123456", "register_type": 2, "password": "123456", "captcha": "508441"}'