FULL=github.com/peter.vaczi/sprinklerd

build:
	go build $(FULL)
	go build $(FULL)/cmd/sprctl

install:
	go install $(FULL)
	go install $(FULL)/cmd/sprctl

test:
	go test $(FULL)
	go test $(FULL)/cmd/sprctl
