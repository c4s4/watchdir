run: clean build
	./watchdir watchdir.yml

build:
	go build watchdir.go

clean:
	rm -f watchdir
