.PHONY: docker_stop docker_start check_docker check_docker_compose run_docker_test docker_test gen_diam test

gen_diam:
	cd diam; ./autogen.sh

test: gen_diam
	go test ./...

##################
# docker endpoints
##################
check_docker: ; @which docker > /dev/null

check_docker_compose: ; @which docker-compose > /dev/null

docker_start:
	@cd docker; docker-compose up -d;

docker_stop:
	@cd docker; docker-compose down

run_docker_test:
	@cd docker; docker-compose exec go-diameter bash -c  "cd go-diameter/ && go test ./..."

docker_test: check_docker check_docker_compose gen_diam docker_start run_docker_test
	$(info )
	$(info --- run `make docker_stop` to stop the test container --- )
