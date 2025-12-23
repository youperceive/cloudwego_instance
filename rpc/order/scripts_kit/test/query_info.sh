#!/bin/bash
kitexcall \
--idl-path ../../idl/order/order.thrift \
--method OrderService/QueryOrderInfo \
--endpoint 127.0.0.1:8002 \
-d '{"id": "694a4341c303adaecb40f8ca"}'

# 694a426acccd3b716b18aa09
# 694a4341c303adaecb40f8ca