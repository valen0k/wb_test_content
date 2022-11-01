.PHONY: build
build:
	sudo docker-compose -f ./docker-compose.yml build
	sudo docker-compose -f ./docker-compose.yml up -d

.PHONY: stop
stop:
	sudo docker stop $$(sudo docker ps -qa)
	sudo docker rm -f $$(sudo docker ps -qa)

.PHONY: re
re: stop build

.PHONY: info
info:
	sudo docker-compose images
	sudo docker-compose ps

.PHONY: logw
logw:
	sudo docker-compose logs -f weather
