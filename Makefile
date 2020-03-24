VERSION := v0.1
IMAGE_NAME := jtgans/feedr

SRCS := Dockerfile feedr.go go.mod go.sum

build: .buildstamp
.buildstamp: $(SRCS)
	docker build -t $(IMAGE_NAME):$(VERSION) .
	docker tag $(IMAGE_NAME):$(VERSION) $(IMAGE_NAME):latest
	touch .buildstamp

push: .pushstamp
.pushstamp: .buildstamp
	docker push $(IMAGE_NAME):$(VERSION)
	docker push $(IMAGE_NAME):latest

clean:
	for image in `docker images |grep $(IMAGE_NAME):$(VERSION) |awk '{ print $3 }'`; do \
		docker rmi $image; \
	done
	for image in `docker images |grep $(IMAGE_NAME):latest |awk '{ print $3 }'`; do \
		docker rmi $image; \
	done

.PHONY: build push run clean
