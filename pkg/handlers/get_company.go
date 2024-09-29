package handlers

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetCompanyHandler struct {
	linkedin_data_pb.GetCompanyServer
	service services.GetCompanyService
	logger  monitor.GRPCLogger
}

func (h *GetCompanyHandler) getCompany(ctx context.Context, in *linkedin_data_pb.GetCompanyRequest) (*linkedin_data_pb.Company, error) {
	company, err := h.service.Exec(ctx, in.GetPublicIdentifier())
	if err != nil {
		if errors.Is(err, dao.ErrCompanyNotFound) {
			return nil, status.Error(codes.NotFound, "company not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to get company: %v", err)
	}

	return &linkedin_data_pb.Company{
		PublicIdentifier:  company.PublicIdentifier,
		Name:              company.Name,
		ProfilePictureUrl: company.ProfilePicture,
	}, nil
}

func (h *GetCompanyHandler) GetCompany(ctx context.Context, in *linkedin_data_pb.GetCompanyRequest) (*linkedin_data_pb.Company, error) {
	res, err := h.getCompany(ctx, in)
	h.logger.Report(ctx, "GetCompany", err)
	return res, err
}

func NewGetCompany(service services.GetCompanyService, logger monitor.GRPCLogger) *GetCompanyHandler {
	return &GetCompanyHandler{
		service: service,
		logger:  logger,
	}
}
