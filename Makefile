bk:
	git add .
	git commit -m update
	git push

redis:
	docker pull redislabs/rebloom:latest
	docker run -d --name redis-bloom -p 6379:6379 redislabs/rebloom:latest
