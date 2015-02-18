VERSION=1.0.0
NAME=watchdir
SOURCE=$(NAME).go
CONFIG=$(NAME).yml
BUILD_DIR=build
DEPLOY=casa@sweetohm.net:/home/web/watchdir
OS_LIST="!plan9"

YELLOW=\033[93m
RED=\033[91m
CLEAR=\033[0m

all: clean test build

clean:
	@echo "${YELLOW}Cleaning generated files${CLEAR}"
	rm -rf $(BUILD_DIR)

test:
	@echo "${YELLOW}Running unit tests${CLEAR}"
	go test

build:
	@echo "${YELLOW}Building application${CLEAR}"
	mkdir -p $(BUILD_DIR)
	go build $(SOURCE)
	mv $(NAME) $(BUILD_DIR)

run: clean test build
	@echo "${YELLOW}Running application${CLEAR}"
	go run $(SOURCE) $(CONFIG)

install: clean test build
	@echo "${YELLOW}Installing application${CLEAR}"
	sudo cp $(BUILD_DIR)/$(NAME) /opt/bin/
	sudo cp $(NAME).init /etc/init.d/watchdir

tag:
	@echo "${YELLOW}Tagging project${CLEAR}"
	git tag "RELEASE-$(VERSION)"
	git push --tag

check:
	@echo "${YELLOW}Chekcing project for release${CLEAR}"
	@if [ `git rev-parse --abbrev-ref HEAD` != "master" ]; then \
		echo "You must release on branch master"; \
		exit 1; \
	fi
	git diff --quiet --exit-code HEAD || (echo "There are uncommitted changes"; exit 1)

binaries: clean test
	@echo "${YELLOW}Generating binaries${CLEAR}"
	mkdir -p $(BUILD_DIR)/$(NAME)
	gox -os=$(OS_LIST) -output=$(BUILD_DIR)/$(NAME)/{{.Dir}}-{{.OS}}-{{.Arch}}

archive: binaries
	@echo "${YELLOW}Generating distribution archive${CLEAR}"
	cp license readme.md $(BUILD_DIR)/$(NAME)
	cd $(BUILD_DIR) && tar cvf $(NAME)-$(VERSION).tar $(NAME)/*
	gzip $(BUILD_DIR)/$(NAME)-$(VERSION).tar

publish: archive
	@echo "${YELLOW}Publishing distribution archive${CLEAR}"
	scp $(BUILD_DIR)/$(NAME)-$(VERSION).tar.gz $(DEPLOY)

release: check publish tag
	@echo "${YELLOW}Application released${CLEAR}"
