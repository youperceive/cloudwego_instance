#!/bin/bash

kitexcall \
--idl-path idl/captcha.thrift \
--method CaptchaService/ValidateCaptcha \
--endpoint 127.0.0.1:8000 \
-f script_kit/test/resp.json