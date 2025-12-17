#!/bin/bash

kitex \
    -module CloudWeGoInstance/VerifyCodeService \
    -service VerifyCodeService \
    -use CloudWeGoInstance/kitex_gen \
    ../../idl/VerifyCodeService/captcha.thrift 