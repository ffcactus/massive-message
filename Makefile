all: receiver sender server

receiver:	
	$(MAKE) -C receiver
	$(MAKE) -C sender
	$(MAKE) -C server
 
clean:
	$(MAKE) -C receiver clean
	$(MAKE) -C sender clean
	$(MAKE) -C server clean

image:
	$(MAKE) -C receiver image
	$(MAKE) -C server image

deploy:
	docker stack deploy -c docker-compose.yml message

undeploy:
	docker stack rm message

.PHONY: receiver sender server all clean image deploy undeploy
