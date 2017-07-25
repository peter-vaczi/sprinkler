FULL=github.com/peter.vaczi/sprinklerd
ADDR=192.168.0.168
OPTS=-s http://$(ADDR):8000

build:
	go build $(FULL)

install:
	go install $(FULL)

test:
	go test $(FULL)

build.arm:
	GOARCH=arm GOOS=linux go build $(FULL)
	scp sprinklerd root@$(ADDR):

setup:
	sprinklerd $(OPTS) device add dev1 --pin 11
	sprinklerd $(OPTS) device add dev2 --pin 12
	sprinklerd $(OPTS) device add dev3 --pin 13
	sprinklerd $(OPTS) device add dev4 --pin 14
	sprinklerd $(OPTS) device add dev5 --pin 15
	sprinklerd $(OPTS) program add pr1
	sprinklerd $(OPTS) program add pr2
	sprinklerd $(OPTS) program adddevice pr1 dev1 --duration 25m
	sprinklerd $(OPTS) program adddevice pr1 dev2 --duration 25m
	sprinklerd $(OPTS) program adddevice pr1 dev3 --duration 10m
	sprinklerd $(OPTS) program adddevice pr1 dev4 --duration 10m
	sprinklerd $(OPTS) program adddevice pr2 dev5 --duration 1h
