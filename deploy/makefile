include .env
export

VERSION ?=0.0.1
API_SERVICE_PORT ?=9011
SERVICE ?=api-server

check:
	sudo docker info

login:
	@echo login
	sudo docker login "${DOCKER_REGISTRY}" -u "${DOCKER_USER}" -p ${DOCKER_PASSWORD}

pull:
	@echo "Pulling latest image for sportNews-$(SERVICE)"
	sudo docker pull hbbbbb.bets8888.com/sport-news/$(SERVICE):$(VERSION)

run:
	make login
	make stop SERVER=$(SERVICE)
	make pull SERVER=$(SERVICE) VERSION=$(VERSION)
	@echo "starting sportNews-$(SERVICE)"
ifeq ($(SERVICE), api-server)
	sudo docker run -d --restart unless-stopped  --name sportNews-$(SERVICE) -p $(API_SERVICE_PORT):9011 -e ASSETS=${ASSETS} -e ASSETS_2=${ASSETS2} \
	-e DB_DATABASE=${DB_DATABASE} -e DB_MASTER_HOST=${DB_HOST} -e DB_MASTER_USERNAME=${DB_USERNAME} -e DB_MASTER_PASSWORD='${DB_PASSWORD}' \
    -e DB_SLAVE_HOST=${DB_HOST} -e DB_SLAVE_USERNAME=${DB_USERNAME} -e DB_SLAVE_PASSWORD='${DB_PASSWORD}' \
    hbbbbb.bets8888.com/sport-news/$(SERVICE):$(VERSION)
else ifeq ($(SERVICE), crawler-server)
	sudo docker run -d --restart unless-stopped --name sportNews-$(SERVICE) \
	-e APP_PROCESS_NEWS=7200 -e APP_PROCESS_RANKING=43200 -e ASSETS=${ASSETS} -e ASSETS_2=${ASSETS2} \
	-e DB_DATABASE=${DB_DATABASE} -e DB_MASTER_HOST=${DB_HOST} -e DB_MASTER_USERNAME=${DB_USERNAME} -e DB_MASTER_PASSWORD='${DB_PASSWORD}' \
	-e DB_SLAVE_HOST=${DB_HOST} -e DB_SLAVE_USERNAME=${DB_USERNAME} -e DB_SLAVE_PASSWORD='${DB_PASSWORD}' \
	hbbbbb.bets8888.com/sport-news/$(SERVICE):$(VERSION)
else
	@echo "Error: Invalid service name $(SERVICE). Please specify a valid service."
endif

stop:
	@echo "Stopping container for sportNews-$(SERVICE)"
	-sudo docker stop sportNews-$(SERVICE)
	-sudo docker rm sportNews-$(SERVICE)