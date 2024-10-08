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

type UpsertUserHandler struct {
	linkedin_data_pb.UpsertUserServer
	service services.UpsertUserService
	logger  monitor.GRPCLogger
}

func (h *UpsertUserHandler) upsertUser(ctx context.Context, in *linkedin_data_pb.UpsertUserRequest) (*linkedin_data_pb.User, error) {
	user, err := h.service.Exec(ctx, in.GetPublicIdentifier(), &models.UpsertUser{
		FirstName:      in.GetFirstName(),
		LastName:       in.GetLastName(),
		ProfilePicture: in.GetProfilePictureBase64(),
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidUpsertUser) || errors.Is(err, dao.ErrInvalidProfilePicture) {
			return nil, status.Errorf(codes.InvalidArgument, "failed to upsert user: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to upsert user: %v", err)
	}

	return &linkedin_data_pb.User{
		PublicIdentifier:  user.PublicIdentifier,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		ProfilePictureUrl: user.ProfilePicture,
	}, nil
}

func (h *UpsertUserHandler) UpsertUser(ctx context.Context, in *linkedin_data_pb.UpsertUserRequest) (*linkedin_data_pb.User, error) {
	res, err := h.upsertUser(ctx, in)
	h.logger.Report(ctx, "UpsertUser", err)
	return res, err
}

func NewUpsertUser(service services.UpsertUserService, logger monitor.GRPCLogger) *UpsertUserHandler {
	return &UpsertUserHandler{
		service: service,
		logger:  logger,
	}
}
