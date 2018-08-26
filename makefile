FULL=github.com/peter-vaczi/sprinkler
ADDR=192.168.0.170
OPTS=-s http://$(ADDR):8000
#RACE=-race

all: build test

build:
	go build -v

install:
	go install -v $(FULL)

test:
	go test $(RACE) -v $(FULL)/core
	go test $(RACE) -v $(FULL)/api

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
	scp sprinkler root@$(ADDR):

test_setup:
	sprinkler $(OPTS) device add dev1 --switch-on-low --pin 9
	sprinkler $(OPTS) device add dev2 --switch-on-low --pin 10
	sprinkler $(OPTS) device add dev3 --switch-on-low --pin 23
	sprinkler $(OPTS) device add dev4 --switch-on-low --pin 24
	sprinkler $(OPTS) device add dev5 --switch-on-low --pin 15
	sprinkler $(OPTS) program add pr1
	sprinkler $(OPTS) program add pr2
	sprinkler $(OPTS) program adddevice pr1 dev1 --duration 5s
	sprinkler $(OPTS) program adddevice pr1 dev2 --duration 5s
	sprinkler $(OPTS) program adddevice pr1 dev3 --duration 3s
	sprinkler $(OPTS) program adddevice pr1 dev4 --duration 3s
	sprinkler $(OPTS) program adddevice pr2 dev5 --duration 10s

test_cleanup:
	-sprinkler $(OPTS) program deldevice pr1 dev1
	-sprinkler $(OPTS) program deldevice pr1 dev2
	-sprinkler $(OPTS) program deldevice pr1 dev3
	-sprinkler $(OPTS) program deldevice pr1 dev4
	-sprinkler $(OPTS) program deldevice pr2 dev5
	-sprinkler $(OPTS) program del pr1
	-sprinkler $(OPTS) program del pr2
	-sprinkler $(OPTS) device del dev1
	-sprinkler $(OPTS) device del dev2
	-sprinkler $(OPTS) device del dev3
	-sprinkler $(OPTS) device del dev4
	-sprinkler $(OPTS) device del dev5
