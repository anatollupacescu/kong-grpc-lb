PROJECT:=$(shell go list -m)

.PHONY: format test run build

format:
	@goimports -w -local $(PROJECT) . ./internal

test:
	@go test -v -trimpath -race -vet all -count=1 -timeout=10s $(shell pwd)/...

run:
	@go run $(shell pwd)

build:
	@docker build -t atlant_api .

.PHONY: compose provision

compose:
	@docker-compose up --build --scale api=2 -d

provision:
	@./configure-load-balancing.sh