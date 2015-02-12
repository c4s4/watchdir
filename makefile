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
