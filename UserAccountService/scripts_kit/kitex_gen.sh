#!/bin/bash

kitex -module CloudWeGoInstance/UserAccountService -service UserAccountService idl/user.thrift
kitex -module CloudWeGoInstance/UserAccountService idl/captcha.thrift