# Trafic Toll Catculator - Microcervices
Work in progress...

## Runu Kafka Container
```
docker compose up -d
```
or
```
docker run -d -p 9092:9092 --name broker apache/kafka:latest

```

## TODO
- [x] Simulate the OBU
- [x] Implement data receiver
- [x] Implement Distancec calculator
- [ ] Implement invoice aggregator
- [ ] Implement GRPC and proto buffers