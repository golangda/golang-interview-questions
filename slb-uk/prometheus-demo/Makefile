.PHONY: run docker-up docker-down

run:
	cd app && go run .

docker-up:
	docker compose -f ops/docker-compose.yml up -d

docker-down:
	docker compose -f ops/docker-compose.yml down