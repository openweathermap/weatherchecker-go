build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/app_linux_amd64
	docker build -t `printenv DOCKER_IMAGE_NAME`:latest .
push: build
	docker push `printenv DOCKER_IMAGE_NAME`:latest
	docker tag `printenv DOCKER_IMAGE_NAME`:latest `printenv DOCKER_IMAGE_NAME`:`cat DOCKER_IMAGE_TAG`
	docker push `printenv DOCKER_IMAGE_NAME`:`cat DOCKER_IMAGE_TAG`