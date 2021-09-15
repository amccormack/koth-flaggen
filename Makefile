.PHONY: all build
all: shell

build-image:
	docker build . -t koth-flaggen

shell: build-image
	docker run -it -v ${PWD}/src/:/src/ koth-flaggen bash