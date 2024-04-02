package linkgrpc

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/internal/database"
	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/internal/link/linkrepository"
	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/internal/user/userrepository"
	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/pkg/pb"
)

type Handler struct {
	pb.UnimplementedLinkServiceServer
	linksRepo linkrepository.Repository
	usersRepo userrepository.Repository
	timeout   time.Duration
}

func New(linksRepo linkrepository.Repository, usersRepo userrepository.Repository, timeout time.Duration) *Handler {
	return &Handler{
		linksRepo: linksRepo,
		usersRepo: usersRepo,
		timeout:   timeout,
	}
}

func (h *Handler) GetLinkByUserID(ctx context.Context, id *pb.GetLinksByUserId) (*pb.ListLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	list, err := h.linksRepo.FindByUserID(ctx, id.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var respList []*pb.Link
	for _, l := range list {
		respList = append(respList, database.LinkToPBLink(&l))
	}

	return &pb.ListLinkResponse{Links: respList}, nil
}

func (h *Handler) CreateLink(ctx context.Context, request *pb.CreateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.validateUserExistence(ctx, request.UserId); err != nil {
		return nil, err
	}

	if _, err := h.linksRepo.Create(ctx, database.PBLinkToLink(request)); err != nil {
		if errors.Is(err, database.ErrConflict) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// implement AMQP publishing here

	return &pb.Empty{}, nil
}

func (h *Handler) GetLink(ctx context.Context, request *pb.GetLinkRequest) (*pb.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	l, err := h.linksRepo.FindByID(ctx, objectID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return database.LinkToPBLink(&l), nil
}

func (h *Handler) UpdateLink(ctx context.Context, request *pb.UpdateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.validateUserExistence(ctx, request.UserId); err != nil {
		return nil, err
	}

	if _, err := h.linksRepo.Update(ctx, database.PBLinkToLink(request)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}

func (h *Handler) DeleteLink(ctx context.Context, request *pb.DeleteLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.linksRepo.Delete(ctx, objectID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}

func (h *Handler) ListLinks(ctx context.Context, request *pb.Empty) (*pb.ListLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	list, err := h.linksRepo.FindAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var respList []*pb.Link
	for _, l := range list {
		respList = append(respList, database.LinkToPBLink(&l))
	}

	return &pb.ListLinkResponse{Links: respList}, nil
}

func (h *Handler) validateUserExistence(ctx context.Context, userID string) error {
	if _, err := primitive.ObjectIDFromHex(userID); err != nil {
		return status.Error(codes.InvalidArgument, "invalid user ID")
	}

	if _, err := h.usersRepo.FindByID(ctx, userID); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return status.Error(codes.NotFound, "user not found")
		}
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
