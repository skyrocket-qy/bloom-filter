.PHONY: run start-redis stop-redis clean

run:
	go run main.go

start-redis:
	docker pull redislabs/rebloom:latest
	docker run -d --name redis-bloom -p 6379:6379 redislabs/rebloom:latest

stop-redis:
	docker stop redis-bloom
	docker rm redis-bloom

clean:
	rm -rf out/
