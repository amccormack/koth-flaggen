.PHONY: all build
all: shell

build-image:
	docker build . -t koth-flaggen

shell: build-image
	docker run -it -v ${PWD}/flaggen/:/flaggen/ koth-flaggen bash

shell-web: build-image
	docker run -it -p 8080:8080 -v ${PWD}/web/:/web/ --workdir /web/ koth-flaggen bash