obu:
	go build -o bin/obu OBU/main.go
	bin/obu

receiver:
	go build -o bin/receiver data-receiver/main.go
	bin/receiver

.PHONY: obu