default:
	@echo "=============Building Production Service============="
	docker build -f ./Dockerfile -t orov.io/${SERVICE_NAME} .

up: 
	@echo "=============Starting Service Locally============="
	docker-compose up -d

build:
	@echo "=============Building Development image============="
	docker-compose build

logs:
	docker-compose logs -f

down:
	@echo "=============Stopping Development Service============="
	docker-compose down

reload: down up

restart: down build up

clean: down
	@echo "=============cleaning up============="
	rm -f api
	docker system prune -f
	docker volume prune -f