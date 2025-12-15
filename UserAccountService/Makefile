all: run

build:
	./build.sh

run: build
	MYSQL_DSN="user_service:user123456@tcp(127.0.0.1:3306)/user_account_db?charset=utf8mb4&parseTime=True&loc=Local" ./output/bootstrap.sh