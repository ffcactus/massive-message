all:
	@echo "Building server ..."
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build

clean:
	@rm -f ./server

image: all
	mv ./server ./build/image/fs/
	docker build -t server.message:latest ./build/image

.PHONY: all clean image
