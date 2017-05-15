FULL=github.com/peter.vaczi/sprinklerd

build:
	go build -i $(FULL)
	go build -i $(FULL)/cmd/sprctl

install:
	go install $(FULL)
	go install $(FULL)/cmd/sprctl

test:
	go test $(FULL)
	go test $(FULL)/cmd/sprctl
