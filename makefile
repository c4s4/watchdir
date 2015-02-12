run: clean test build
	./watchdir watchdir.yml


test:
	go test

build:
	go build watchdir.go

clean:
	rm -f watchdir
