package update

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"yadro.com/course/api/core"
	updatepb "yadro.com/course/proto/update"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (m *MockClient) Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*updatepb.StatusReply, error) {
	return &updatepb.StatusReply{
		Status: 1,
	}, nil
}

func (m *MockClient) Stats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*updatepb.StatsReply, error) {
	return &updatepb.StatsReply{
		WordsTotal:    1,
		WordsUnique:   1,
		ComicsTotal:   1,
		ComicsFetched: 1,
	}, nil
}

func (m *MockClient) Update(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (m *MockClient) Drop(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

type DummyLogger struct{}

func (d DummyLogger) Error(msg string, keysAndValues ...interface{}) {}

func TestPingOK(t *testing.T) {

	c := Client{
		log:    DummyLogger{},
		client: &MockClient{},
	}
	err := c.Ping(context.Background())
	require.NoError(t, err)
}

func TestStatusOK(t *testing.T) {

	c := Client{
		log:    DummyLogger{},
		client: &MockClient{},
	}
	_, err := c.Status(context.Background())
	require.NoError(t, err)
}

func TestStatsOK(t *testing.T) {

	c := Client{
		log:    DummyLogger{},
		client: &MockClient{},
	}

	good_stats := core.UpdateStats{
		WordsTotal:    1,
		WordsUnique:   1,
		ComicsTotal:   1,
		ComicsFetched: 1,
	}
	stats, err := c.Stats(context.Background())
	require.NoError(t, err)
	require.Equal(t, stats, good_stats)
}

func TestUpdateOK(t *testing.T) {

	c := Client{
		log:    DummyLogger{},
		client: &MockClient{},
	}
	err := c.Update(context.Background())
	require.NoError(t, err)
}

func TestDropOK(t *testing.T) {

	c := Client{
		log:    DummyLogger{},
		client: &MockClient{},
	}
	err := c.Drop(context.Background())
	require.NoError(t, err)
}
