package handlers

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetCompanyLastUpdateHandler struct {
	linkedin_data_pb.GetCompanyLastUpdateServer
	service services.GetCompanyLastUpdateService
	logger  monitor.GRPCLogger
}

func (h *GetCompanyLastUpdateHandler) getCompanyLastUpdate(ctx context.Context, in *linkedin_data_pb.GetCompanyLastUpdateRequest) (*linkedin_data_pb.GetCompanyLastUpdateResponse, error) {
	companyLastUpdate, err := h.service.Exec(ctx, in.GetPublicIdentifier())
	if err != nil {
		if errors.Is(err, dao.ErrCompanyNotFound) {
			return nil, status.Error(codes.NotFound, "company not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to get company last update: %v", err)
	}

	return &linkedin_data_pb.GetCompanyLastUpdateResponse{UpdatedAt: timestamppb.New(lo.FromPtr(companyLastUpdate))}, nil
}

func (h *GetCompanyLastUpdateHandler) GetCompanyLastUpdate(ctx context.Context, in *linkedin_data_pb.GetCompanyLastUpdateRequest) (*linkedin_data_pb.GetCompanyLastUpdateResponse, error) {
	res, err := h.getCompanyLastUpdate(ctx, in)
	h.logger.Report(ctx, "GetCompanyLastUpdate", err)
	return res, err
}

func NewGetCompanyLastUpdate(service services.GetCompanyLastUpdateService, logger monitor.GRPCLogger) *GetCompanyLastUpdateHandler {
	return &GetCompanyLastUpdateHandler{
		service: service,
		logger:  logger,
	}
}
