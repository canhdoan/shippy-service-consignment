// shippy-service-consignment/handler.go
package main

import (
	"context"
	pb "github.com/canhdoan/shippy-service-consignment/proto/consignment"
	vesselProto "github.com/canhdoan/shippy-service-vessel/proto/vessel"
	"gopkg.in/mgo.v2"
)

type service struct {
	session      *mgo.Session
	vesselClient vesselProto.VesselServiceClient
}

func (s *service) GetRepo() Repository {
	return &ConsignmentRepository{s.session.Clone()}
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	// Clone a db session and close it after done
	repo := s.GetRepo()
	defer repo.Close()

	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(),
		&vesselProto.Specification{
			MaxWeight: req.Weight,
			Capacity:  int32(len(req.Containers)),
		})

	if err != nil {
		return err
	}

	// set the vailable vessel for new consignment
	req.VesselId = vesselResponse.Vessel.Id

	// save
	err = repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = req

	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	// Clone a db session and close it after done
	repo := s.GetRepo()
	defer repo.Close()

	consignments, err := repo.GetAll()
	if err != nil {
		return err
	}
	res.Consignments = consignments
	return nil
}
