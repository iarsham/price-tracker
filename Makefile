down:
	docker compose down

up:
	docker compose up --build -d

log-bot:
	docker compose logs -f price-tracker-bot