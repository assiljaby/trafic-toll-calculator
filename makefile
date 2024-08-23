obu:
	go build -o bin/obu OBU/main.go
	bin/obu

receiver:
	go build -o bin/receiver ./data-receiver
	bin/receiver

calculator:
	go build -o bin/calculator ./distance-calculator
	bin/calculator

agg:
	go build -o bin/aggregator ./aggregator
	bin/aggregator

invoicer:
	go build -o bin/invoicer ./invoicer
	bin/invoicer

proto:
	protoc --go_out=. --go_opt=paths=source_relative types/ptypes.proto

.PHONY: obu invoicer