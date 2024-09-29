package handlers

import (
	"context"
	"github.com/in-rich/lib-go/monitor"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListUsersHandler struct {
	linkedin_data_pb.ListUsersServer
	service services.ListUsersService
	logger  monitor.GRPCLogger
}

func (h *ListUsersHandler) listUsers(ctx context.Context, in *linkedin_data_pb.ListUsersRequest) (*linkedin_data_pb.ListUsersResponse, error) {
	users, err := h.service.Exec(ctx, in.GetPublicIdentifiers())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	res := &linkedin_data_pb.ListUsersResponse{
		Users: make([]*linkedin_data_pb.User, len(users)),
	}
	for i, user := range users {
		res.Users[i] = &linkedin_data_pb.User{
			PublicIdentifier:  user.PublicIdentifier,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			ProfilePictureUrl: user.ProfilePicture,
		}
	}

	return res, nil
}

func (h *ListUsersHandler) ListUsers(ctx context.Context, in *linkedin_data_pb.ListUsersRequest) (*linkedin_data_pb.ListUsersResponse, error) {
	res, err := h.listUsers(ctx, in)
	h.logger.Report(ctx, "ListUsers", err)
	return res, err
}

func NewListUsers(service services.ListUsersService, logger monitor.GRPCLogger) *ListUsersHandler {
	return &ListUsersHandler{
		service: service,
		logger:  logger,
	}
}
