# MONGODB_URI = mongodb://root:root123456@localhost:27017/order_db?authSource=admin

all: run

build:
	./build.sh

run: build
	MONGODB_URI=$(MONGODB_URI) ./output/bootstrap.sh