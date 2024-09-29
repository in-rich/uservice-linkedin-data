package handlers

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/models"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpsertCompanyHandler struct {
	linkedin_data_pb.UpsertCompanyServer
	service services.UpsertCompanyService
	logger  monitor.GRPCLogger
}

func (h *UpsertCompanyHandler) upsertCompany(ctx context.Context, in *linkedin_data_pb.UpsertCompanyRequest) (*linkedin_data_pb.Company, error) {
	company, err := h.service.Exec(ctx, in.GetPublicIdentifier(), &models.UpsertCompany{
		Name:           in.GetName(),
		ProfilePicture: in.GetProfilePictureBase64(),
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidUpsertCompany) || errors.Is(err, dao.ErrInvalidProfilePicture) {
			return nil, status.Errorf(codes.InvalidArgument, "failed to upsert company: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to upsert company: %v", err)
	}

	return &linkedin_data_pb.Company{
		PublicIdentifier:  company.PublicIdentifier,
		Name:              company.Name,
		ProfilePictureUrl: company.ProfilePicture,
	}, nil
}

func (h *UpsertCompanyHandler) UpsertCompany(ctx context.Context, in *linkedin_data_pb.UpsertCompanyRequest) (*linkedin_data_pb.Company, error) {
	res, err := h.upsertCompany(ctx, in)
	h.logger.Report(ctx, "UpsertCompany", err)
	return res, err
}

func NewUpsertCompany(service services.UpsertCompanyService, logger monitor.GRPCLogger) *UpsertCompanyHandler {
	return &UpsertCompanyHandler{
		service: service,
		logger:  logger,
	}
}
