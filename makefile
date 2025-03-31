include .env
export

VERSION ?= 0.0.1
API_SERVICE_PORT ?=9011
SERVICE ?=api-server

check:
	docker info

login:
	@echo login
	docker login "${DOCKER_REGISTRY}" -u "${DOCKER_USER}" -p ${DOCKER_PASSWORD}

build:
	@echo build hbbbbb.bets8888.com/sport-news/$(SERVICE):$(VERSION)
ifeq ($(SERVICE), api-server)
	docker build -t hbbbbb.bets8888.com/sport-news/$(SERVICE):$(VERSION) -f .\docker\app\Dockerfile .
else ifeq ($(SERVICE), crawler-server)
	docker build -t hbbbbb.bets8888.com/sport-news/$(SERVICE):$(VERSION) -f .\docker\crawler\Dockerfile .
else
endif


push:
	@echo push hbbbbb.bets8888.com/sport-news/$(SERVICE):$(VERSION)
	make login
	docker push hbbbbb.bets8888.com/sport-news/$(SERVICE):$(VERSION)

run:
	make login
	make stop SERVER=$(SERVICE)
	@echo start sportNews-$(SERVICE)
ifeq ($(SERVICE), api-server)
	docker run --rm -d --name sportNews-$(SERVICE) -p $(API_SERVER_PORT):9011 -e DB_MASTER_HOST=${DB_HOST} -e DB_MASTER_USERNAME=${DB_USERNAME} -e DB_MASTER_PASSWORD=${DB_PASSWORD} -e DB_SLAVE_HOST=${DB_HOST} -e DB_SLAVE_USERNAME=${DB_USERNAME} -e DB_SLAVE_PASSWORD=${DB_PASSWORD} hbbbbb.bets8888.com/sport-news/$(SERVICE):$(VERSION)
else ifeq ($(SERVICE), crawler-server)
	docker run --rm -d --name sportNews-$(SERVICE) -e DB_MASTER_HOST=${DB_HOST} -e DB_MASTER_USERNAME=${DB_USERNAME} -e DB_MASTER_PASSWORD=${DB_PASSWORD} -e DB_SLAVE_HOST=${DB_HOST} -e DB_SLAVE_USERNAME=${DB_USERNAME} -e DB_SLAVE_PASSWORD=${DB_PASSWORD} hbbbbb.bets8888.com/sport-news/$(SERVICE):$(VERSION)
else
endif

stop:
	@echo stop sportNews-$(SERVICE)
	-docker stop sportNews-$(SERVICE)
	-docker rm sportNews-$(SERVICE)
