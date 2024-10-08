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

type GetUserHandler struct {
	linkedin_data_pb.GetUserServer
	service services.GetUserService
	logger  monitor.GRPCLogger
}

func (h *GetUserHandler) getUser(ctx context.Context, in *linkedin_data_pb.GetUserRequest) (*linkedin_data_pb.User, error) {
	user, err := h.service.Exec(ctx, in.GetPublicIdentifier())
	if err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &linkedin_data_pb.User{
		PublicIdentifier:  user.PublicIdentifier,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		ProfilePictureUrl: user.ProfilePicture,
	}, nil
}

func (h *GetUserHandler) GetUser(ctx context.Context, in *linkedin_data_pb.GetUserRequest) (*linkedin_data_pb.User, error) {
	res, err := h.getUser(ctx, in)
	h.logger.Report(ctx, "GetUser", err)
	return res, err
}

func NewGetUser(service services.GetUserService, logger monitor.GRPCLogger) *GetUserHandler {
	return &GetUserHandler{
		service: service,
		logger:  logger,
	}
}
