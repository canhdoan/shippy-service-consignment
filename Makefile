build:
	protoc -I. --go_out==plugins=micro:. proto/consignment/consignment.proto
	docker build -t shippy-service-consignment .

run:
	docker run -p 5100:5100 -e MICRO_SERVER_ADDRESS=:5100 shippy-service-consignment