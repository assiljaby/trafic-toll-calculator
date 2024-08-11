obu:
	go build -o bin/obu OBU/main.go
	bin/obu

receiver:
	go build -o bin/receiver ./data-receiver
	bin/receiver

calculator:
	go build -o bin/calculator ./distance-calculator
	bin/calculator

.PHONY: obu