// shippy-service-consignment/main.go
package main

import (
	"context"
	"flag"
	"fmt"

	// Import the generated protobuf code
	"github.com/micro/go-micro"
	pb "gitlab.asoft-python.com/g-canhdoan/golang-training/shippy-service-consignment/proto/consignment"
)

var port = flag.String("l", ":5100", "Specify the port that the server will listen on")

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type Repository struct {
	consignments []*pb.Consignment
}

// Create a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo repository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetConsignmentRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}

	// create a new service
	srv := micro.NewService(
		micro.Name("shippy.service.consignment"),
	)
	srv.Init()

	// Register our service with the gRPC server
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
