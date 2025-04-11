package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	updatepb "yadro.com/course/proto/update"
	"yadro.com/course/update/core"
)

func NewServer(service core.Updater) *Server {
	return &Server{service: service}
}

type Server struct {
	updatepb.UnimplementedUpdateServer
	service core.Updater
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Server) Status(ctx context.Context, _ *emptypb.Empty) (*updatepb.StatusReply, error) {
	res := s.service.Status(ctx)

	switch res {
	case core.StatusIdle:
		return &updatepb.StatusReply{Status: updatepb.Status_STATUS_IDLE}, nil
	case core.StatusRunning:
		return &updatepb.StatusReply{Status: updatepb.Status_STATUS_RUNNING}, nil
	}
	return nil, status.Error(codes.Internal, "Unknown status from service")

}

func (s *Server) Update(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.service.Update(ctx); err != nil {
		if errors.Is(err, core.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "Update is running")
		}
		return nil, err
	}
	return nil, nil
}

func (s *Server) Stats(ctx context.Context, _ *emptypb.Empty) (*updatepb.StatsReply, error) {
	res, err := s.service.Stats(ctx)
	if err != nil {
		return nil, err
	}
	return &updatepb.StatsReply{
		WordsTotal:    int64(res.WordsTotal),
		WordsUnique:   int64(res.WordsUnique),
		ComicsFetched: int64(res.ComicsFetched),
		ComicsTotal:   int64(res.ComicsTotal),
	}, nil
}

func (s *Server) Drop(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.service.Drop(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}
