package handlers

import (
	"context"
	"errors"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetUserLastUpdateHandler struct {
	linkedin_data_pb.GetUserLastUpdateServer
	service services.GetUserLastUpdateService
}

func (h *GetUserLastUpdateHandler) GetUserLastUpdate(ctx context.Context, in *linkedin_data_pb.GetUserLastUpdateRequest) (*linkedin_data_pb.GetUserLastUpdateResponse, error) {
	userLastUpdate, err := h.service.Exec(ctx, in.GetPublicIdentifier())
	if err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to get user last update: %v", err)
	}

	return &linkedin_data_pb.GetUserLastUpdateResponse{UpdatedAt: timestamppb.New(lo.FromPtr(userLastUpdate))}, nil
}

func NewGetUserLastUpdate(service services.GetUserLastUpdateService) *GetUserLastUpdateHandler {
	return &GetUserLastUpdateHandler{
		service: service,
	}
}
