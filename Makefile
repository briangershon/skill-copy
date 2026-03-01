VERSION := $(shell git describe --tags --always --dirty)

build:
	go build -ldflags "-X main.version=$(VERSION)" -o skill-copy .

install:
	go install -ldflags "-X main.version=$(VERSION)" .

tag:
	@test -n "$(TAG)" || (echo "Usage: make tag TAG=v1.0.0" && exit 1)
	git tag $(TAG)
	git push origin $(TAG)

clean:
	rm -f skill-copy
