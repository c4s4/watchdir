VERSION=1.0.0

all: clean test build

clean:
	rm -f watchdir

test:
	go test

build:
	go build watchdir.go

run: clean test build
	go run watchdir.go watchdir.yml

install: clean test build
	sudo cp watchdir /opt/bin/
	sudo cp watchdir.init /etc/init.d/watchdir

release: clean test build
	@if [ `git rev-parse --abbrev-ref HEAD` != "master" ]; then \
		echo "You must release on branch master"; \
		exit 1; \
	fi
	git diff --quiet --exit-code HEAD || (echo "There are uncommitted changes"; exit 1)
	git tag "RELEASE-$(VERSION)"
	git push --tag
