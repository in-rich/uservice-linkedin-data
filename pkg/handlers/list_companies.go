package handlers

import (
	"context"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListCompaniesHandler struct {
	linkedin_data_pb.ListCompaniesServer
	service services.ListCompaniesService
}

func (h *ListCompaniesHandler) ListCompanies(ctx context.Context, in *linkedin_data_pb.ListCompaniesRequest) (*linkedin_data_pb.ListCompaniesResponse, error) {
	companies, err := h.service.Exec(ctx, in.GetPublicIdentifiers())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list companies: %v", err)
	}

	res := &linkedin_data_pb.ListCompaniesResponse{
		Companies: make([]*linkedin_data_pb.Company, len(companies)),
	}
	for i, company := range companies {
		res.Companies[i] = &linkedin_data_pb.Company{
			PublicIdentifier:  company.PublicIdentifier,
			Name:              company.Name,
			ProfilePictureUrl: company.ProfilePicture,
		}
	}

	return res, nil

}

func NewListCompanies(service services.ListCompaniesService) *ListCompaniesHandler {
	return &ListCompaniesHandler{
		service: service,
	}
}
