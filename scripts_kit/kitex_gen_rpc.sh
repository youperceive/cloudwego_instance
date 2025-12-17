#!/bin/bash

kitex \
    -module CloudWeGoInstance \
    -service VerifyCodeService \
    -use CloudWeGoInstance/kitex_gen \
    -gen-path rpc/VerifyCodeService/ \
    idl/VerifyCodeService/captcha.thrift 