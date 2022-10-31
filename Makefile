.PHONY: run
run:
	docker-compose -f ./deploy/docker-compose.yml up

.PHONY: stop
stop:
	docker-compose -f ./deploy/docker-compose.yml down

.PHONY: stop-clean-volumes
stop-clean-volumes:
	docker-compose -f ./deploy/docker-compose.yml down -v

.PHONY: restart
restart:
	docker-compose -f ./deploy/docker-compose.yml down
	docker image rm -f deploy-asker
	docker-compose -f ./deploy/docker-compose.yml up

.PHONY: test-local
test:
	go test ./test -v -url=http://localhost

# Need to install: https://github.com/go-swagger/go-swagger
.PHONY: gen-swagger
gen-swagger:
	GO111MODULE=off swagger generate spec -o ./api/swagger.yaml --scan-models