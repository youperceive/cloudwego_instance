#!/bin/bash
kitexcall \
--idl-path ../../idl/order/order.thrift \
--method OrderService/Update \
--endpoint 127.0.0.1:8002 \
-d '{"id": "694a4341c303adaecb40f8ca", "status": 3}'