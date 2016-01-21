clean:
	rm -rf ./bin ./bindata
	rm -rf ./data/ui/bundle/*
generate: clean
	go generate
build: generate
	GOOS=linux GOARCH=amd64 go build -o ./bin/app_linux_amd64
	docker build -t `printenv DOCKER_IMAGE_NAME`:dev .
push:
	docker push `printenv DOCKER_IMAGE_NAME`:dev
push-dev: push
	docker tag -f `printenv DOCKER_IMAGE_NAME`:dev `printenv DOCKER_IMAGE_NAME`:`cat DOCKER_IMAGE_TAG`
	docker push `printenv DOCKER_IMAGE_NAME`:`cat DOCKER_IMAGE_TAG`
push-stage: push
	docker tag -f `printenv DOCKER_IMAGE_NAME`:dev `printenv DOCKER_IMAGE_NAME`:stage
	docker push `printenv DOCKER_IMAGE_NAME`:stage
push-stable: push
	docker tag -f `printenv DOCKER_IMAGE_NAME`:dev `printenv DOCKER_IMAGE_NAME`:`cat DOCKER_IMAGE_TAG`
	docker tag -f `printenv DOCKER_IMAGE_NAME`:`cat DOCKER_IMAGE_TAG` `printenv DOCKER_IMAGE_NAME`:latest
	docker push `printenv DOCKER_IMAGE_NAME`:`cat DOCKER_IMAGE_TAG`
	docker push `printenv DOCKER_IMAGE_NAME`:latest
