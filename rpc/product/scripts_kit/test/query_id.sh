#!/bin/bash
kitexcall \
--idl-path ../../idl/order/order.thrift \
--method OrderService/QueryOrderId \
--endpoint 127.0.0.1:8002 \
-d '{"type": 1, "user_id": 100001}'

# 694a426acccd3b716b18aa09
# 694a4341c303adaecb40f8ca