all:
	@echo "Building notification ..."
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build

clean:
	@rm -f ./notification

image: all
	mv ./notification ./build/image/fs/
	docker build -t notification.message:latest ./build/image

.PHONY: all clean image
