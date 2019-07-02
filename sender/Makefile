all:
	@echo "Building sender ..."
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build

clean:
	@rm -f ./sender

.PHONY: all clean
