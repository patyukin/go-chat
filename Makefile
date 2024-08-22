.PHONY: upb down upv downv bup r restart

upb:
	docker compose up -d --build

downv:
	docker compose down -v --remove-orphans

up:
	docker compose up -d

down:
	docker compose down

restart:
	make downv
	make upb

r:
	docker compose down
	docker compose up -d
