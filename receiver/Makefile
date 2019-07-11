all:
	@echo "Building receiver ..."
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build

clean:
	@rm -f ./receiver

image: all
	mv ./receiver ./build/image/fs/
	docker build -t receiver.message:latest ./build/image

.PHONY: all clean image
