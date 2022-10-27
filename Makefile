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
