verify_code_service_addr="127.0.0.1:8000"
user_account_service_addr="127.0.0.1:8001"
order_service_addr="127.0.0.1:8002"
MYSQL_DSN="root:123456@tcp(127.0.0.1:3306)/product?charset=utf8mb4&parseTime=True"

all: run

build:
	./build.sh

run: build
	verify_code_service_addr=$(verify_code_service_addr) \
	user_account_service_addr=$(user_account_service_addr) \
	order_service_addr=$(order_service_addr) \
	MYSQL_DSN=$(MYSQL_DSN) \
	./output/bootstrap.sh
