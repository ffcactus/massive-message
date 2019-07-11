all: receiver sender server notification

receiver:	
	$(MAKE) -C receiver
sender:
	$(MAKE) -C sender
server:
	$(MAKE) -C server
notification:
	$(MAKE) -C notification
 
clean:
	$(MAKE) -C receiver clean
	$(MAKE) -C sender clean
	$(MAKE) -C server clean
	$(MAKE) -C notification clean
image:
	$(MAKE) -C receiver image
	$(MAKE) -C server image
	$(MAKE) -C notification image
	$(MAKE) -C postgres image
	$(MAKE) -C rabbitmq image
	$(MAKE) -C httpd image

deploy:
	docker stack deploy -c docker-compose.yml message

undeploy:
	docker stack rm message

.PHONY: receiver sender server notification all clean image deploy undeploy
