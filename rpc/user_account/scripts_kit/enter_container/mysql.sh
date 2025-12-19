#!/bin/bash

MYSQL_PWD="user123456"
# 容器名
CONTAINER_NAME="user_account-mysql-1"
# 数据库用户名
USER="user_service"
# 数据库名
DB_NAME="user_account_db"

核心命令：非交互执行SQL
docker exec -i $CONTAINER_NAME mysql -u $USER -p$MYSQL_PWD -t -e "
USE $DB_NAME;
SELECT * FROM user;
"

# docker exec -it $CONTAINER_NAME mysql -u $USER -p$MYSQL_PWD