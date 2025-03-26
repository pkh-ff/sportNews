include .env
export

VERSION ?= 0.0.1
API_SERVER_PORT ?=9011

check:
	docker info

login:
	docker login "${DOCKER_REGISTRY}" -u "${DOCKER_USER}" -p ${DOCKER_PASSWORD}

build-api-server:
	docker build -t hbbbbb.bets8888.com/sport-news/api-server:$(VERSION) -f .\docker\app\Dockerfile .

push-api-server:
	docker push hbbbbb.bets8888.com/sport-news/api-server:$(VERSION)

run-api-server:
	docker run --rm -p $(API_SERVER_PORT):9011 -e DB_MASTER_HOST=${DB_HOST} -e DB_MASTER_USERNAME=${DB_USERNAME} -e DB_MASTER_PASSWORD=${DB_PASSWORD} -e DB_SLAVE_HOST=${DB_HOST} -e DB_SLAVE_USERNAME=${DB_USERNAME} -e DB_SLAVE_PASSWORD=${DB_PASSWORD} hbbbbb.bets8888.com/sport-news/api-server:$(VERSION)