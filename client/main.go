// shippy-cli-consignment/main.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"context"

	micro "github.com/micro/go-micro"
	pb "gitlab.asoft-python.com/g-canhdoan/golang-training/shippy-service-consignment/proto/consignment"
)

const (
	address         = "localhost:5100"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	// // connect to service
	// conn, err := grpc.Dial(address, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("Did not connect: %v", err)
	// }
	// defer conn.Close()
	// client := pb.NewShippingServiceClient(conn)

	service := micro.NewService(micro.Name("shippy.consignment.cli"))
	service.Init()

	client := pb.NewShippingServiceClient("shippy.consignment.service", service.Client())

	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetConsignmentRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
