FULL=github.com/peter.vaczi/sprinklerd

build:
	go build $(FULL)

install:
	go install $(FULL)

test:
	go test $(FULL)

build.arm:
	GOARCH=arm GOOS=linux go build $(FULL)

