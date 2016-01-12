clean:
	rm -rf ./bin ./bindata
build: clean
	go generate
	GOOS=linux GOARCH=amd64 go build -o ./bin/app_linux_amd64
	docker build -t `printenv DOCKER_IMAGE_NAME`:latest .
push: build
	docker push `printenv DOCKER_IMAGE_NAME`:latest
push-release: push
	docker tag -f `printenv DOCKER_IMAGE_NAME`:latest `printenv DOCKER_IMAGE_NAME`:`cat DOCKER_IMAGE_TAG`
	docker push `printenv DOCKER_IMAGE_NAME`:`cat DOCKER_IMAGE_TAG`
push-stage: push
	docker tag -f `printenv DOCKER_IMAGE_NAME`:latest `printenv DOCKER_IMAGE_NAME`:stage
	docker push `printenv DOCKER_IMAGE_NAME`:stage
