#.PHONY: all
#all:
#	@go run app/cmd/weatherserver/main.go

.PHONY: build
build:
	sudo docker-compose -f ./docker-compose.yml build
	sudo docker-compose -f ./docker-compose.yml up -d

.PHONY: stop
stop:
	sudo docker stop $$(sudo docker ps -qa)
	sudo docker rm -f $$(sudo docker ps -qa)
