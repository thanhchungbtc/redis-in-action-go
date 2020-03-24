redis-server:
	docker-compose up -d

redis-cli:
	docker exec -it redis redis-cli
