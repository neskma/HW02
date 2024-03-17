package linkgrpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
)

var _ pb.LinkServiceServer = (*Handler)(nil)

func New(linksRepository linksRepository, timeout time.Duration) *Handler {
	return &Handler{linksRepository: linksRepository, timeout: timeout}
}

type Handler struct {
	pb.UnimplementedLinkServiceServer
	linksRepository linksRepository
	timeout         time.Duration
}

func (h Handler) CreateLink(ctx context.Context, request *pb.CreateLinkRequest) (*emptypb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// TODO implement me
	return nil, status.Error(codes.Unimplemented, codes.Unimplemented.String())
}

func (h Handler) GetLink(ctx context.Context, request *pb.GetLinkRequest) (*pb.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// TODO implement me
	return nil, status.Error(codes.Unimplemented, codes.Unimplemented.String())
}

func (h Handler) UpdateLink(ctx context.Context, request *pb.UpdateLinkRequest) (*emptypb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// TODO implement me
	return nil, status.Error(codes.Unimplemented, codes.Unimplemented.String())
}

func (h Handler) DeleteLink(ctx context.Context, request *pb.DeleteLinkRequest) (*emptypb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// TODO implement me
	return nil, status.Error(codes.Unimplemented, codes.Unimplemented.String())
}

func (h Handler) ListLinks(ctx context.Context, request *emptypb.Empty) (*pb.ListLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// TODO implement me
	return nil, status.Error(codes.Unimplemented, codes.Unimplemented.String())
}
