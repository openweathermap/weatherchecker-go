clean:
	rm -rf ./bin ./bindata
	rm -rf ./data/ui/bundle/*
generate: clean
	go generate
compile: generate
	GOOS=linux GOARCH=amd64 go build -o ./bin/app_linux_amd64

dcbuild: compile
	docker-compose build
dcrun: dcbuild
	docker-compose up -d

build: compile
	docker build -t `printenv DOCKER_IMAGE_NAME`-backend:dev .
	docker build -t `printenv DOCKER_IMAGE_NAME`-frontend:dev ./.docker/nginx
push-stage:
	docker tag -f `printenv DOCKER_IMAGE_NAME`-backend:dev `printenv DOCKER_IMAGE_NAME`-backend:stage
	docker tag -f `printenv DOCKER_IMAGE_NAME`-frontend:dev `printenv DOCKER_IMAGE_NAME`-frontend:stage
	docker push `printenv DOCKER_IMAGE_NAME`-backend:stage
	docker push `printenv DOCKER_IMAGE_NAME`-frontend:stage
push-stable:
	docker tag -f `printenv DOCKER_IMAGE_NAME`-backend:dev `printenv DOCKER_IMAGE_NAME`-backend:latest
	docker tag -f `printenv DOCKER_IMAGE_NAME`-frontend:dev `printenv DOCKER_IMAGE_NAME`-frontend:latest
	docker push `printenv DOCKER_IMAGE_NAME`-backend:latest
	docker push `printenv DOCKER_IMAGE_NAME`-frontend:latest
