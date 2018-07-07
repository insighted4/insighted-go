PROJ=insighted-go
ORG_PATH=github.com/insighted4
REPO_PATH=$(ORG_PATH)/$(PROJ)

all: test

deps:
	go get -d -v $(REPO_PATH)/...

updatedeps:
	go get -d -v -u -f $(REPO_PATH)/...

testdeps:
	go get -d -v -t $(REPO_PATH)/...

updatetestdeps:
	go get -d -v -t -u -f $(REPO_PATH)/...

build: deps
	go build $(REPO_PATH)/...

install: deps
	go install $(REPO_PATH)/...

lint: testdeps
	go get -v github.com/golang/lint/golint
	for file in $$(find . -name '*.go' | grep -v '\.pb\.go\|\.pb\.gw\.go\|examples\|pubsub\/aws\/awssub_test\.go' | grep -v 'server\/kit\/kitserver_pb_test\.go'); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

vet: testdeps
	go vet $(REPO_PATH)/...

errcheck: testdeps
	go get -v github.com/kisielk/errcheck
	errcheck -ignoretests $(REPO_PATH)/...

pretest: lint vet # errcheck

test: testdeps pretest
	go test $(REPO_PATH)/...

clean:
	go clean -i $(REPO_PATH)/...

coverage: testdeps
	./coverage.sh --coveralls

.PHONY: \
	all \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	build \
	install \
	lint \
	vet \
	errcheck \
	pretest \
	test \
	clean \
	coverage
