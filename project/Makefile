FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerApp
MAIL_BINARY=mailApp
LISTENER_BINARY=listenerApp
PAYMENT_BINARY=paymentApp
ENQUIRY_BINARY=enquiryApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mail build_listener build_payment
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && go build -tags dynamic -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_logger: builds the broker binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

## build_broker: builds the auth binary as a linux executable
build_auth:
	@echo "Building authApp binary..."
	cd ../authentication-service && go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd ../mail-service && go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Done!"

## build_listener: builds the mail binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ../listener-service && go build -o ${LISTENER_BINARY} ./cmd
	@echo "Done!"

## build_payment: builds the mail binary as a linux executable
build_payment:
	@echo "Building payment binary..."
	cd ../payment-service && go build -o ${PAYMENT_BINARY} ./cmd/api
	@echo "Done!"

## build_enquiry: builds the mail binary as a linux executable
build_enquiry:
	@echo "Building enquiry binary..."
	cd ../enquiry-service && go build -o ${ENQUIRY_BINARY} ./cmd/api
	@echo "Done!"



## build_front: builds the front end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"