FULL=github.com/peter-vaczi/sprinklerd
ADDR=192.168.0.168
OPTS=-s http://$(ADDR):8000

all: build test

build:
	go build -v

install:
	go install -v $(FULL)

test:
	go test -v $(FULL)/core
	go test -v $(FULL)/api

cover:
	go test -cover -coverprofile cover.core.out $(FULL)/core
	go tool cover -func cover.core.out
	go tool cover -html cover.core.out -o cover.core.html
	go test -cover -coverprofile cover.api.out $(FULL)/api
	go tool cover -func cover.api.out
	go tool cover -html cover.api.out -o cover.api.html

vendorinstall:
	glide install

build.arm:
	GOARCH=arm GOOS=linux go build $(FULL)
	scp sprinklerd root@$(ADDR):

test_setup:
	sprinklerd $(OPTS) device add dev1 --switch-on-low --pin 9
	sprinklerd $(OPTS) device add dev2 --switch-on-low --pin 10
	sprinklerd $(OPTS) device add dev3 --switch-on-low --pin 23
	sprinklerd $(OPTS) device add dev4 --switch-on-low --pin 24
	sprinklerd $(OPTS) device add dev5 --switch-on-low --pin 15
	sprinklerd $(OPTS) program add pr1
	sprinklerd $(OPTS) program add pr2
	sprinklerd $(OPTS) program adddevice pr1 dev1 --duration 5s
	sprinklerd $(OPTS) program adddevice pr1 dev2 --duration 5s
	sprinklerd $(OPTS) program adddevice pr1 dev3 --duration 3s
	sprinklerd $(OPTS) program adddevice pr1 dev4 --duration 3s
	sprinklerd $(OPTS) program adddevice pr2 dev5 --duration 10s

test_cleanup:
	sprinklerd $(OPTS) program deldevice pr1 dev1
	sprinklerd $(OPTS) program deldevice pr1 dev2
	sprinklerd $(OPTS) program deldevice pr1 dev3
	sprinklerd $(OPTS) program deldevice pr1 dev4
	sprinklerd $(OPTS) program deldevice pr2 dev5
	sprinklerd $(OPTS) program del pr1
	sprinklerd $(OPTS) program del pr2
	sprinklerd $(OPTS) device del dev1
	sprinklerd $(OPTS) device del dev2
	sprinklerd $(OPTS) device del dev3
	sprinklerd $(OPTS) device del dev4
	sprinklerd $(OPTS) device del dev5
