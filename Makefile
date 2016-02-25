clean:
	rm -rf ./bin
	rm -rf ./ui/bundle/* ./.docker/nginx/fs/etc/nginx/html
generate: clean
	go generate
	cp -r ./ui ./.docker/nginx/fs/etc/nginx/html
compile: generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/app_linux_amd64

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
