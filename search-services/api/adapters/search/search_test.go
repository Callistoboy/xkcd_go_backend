package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	searchpb "yadro.com/course/proto/search"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (m *MockClient) Search(ctx context.Context, phrase *searchpb.SearchRequest, opts ...grpc.CallOption) (*searchpb.SearchReply, error) {
	return &searchpb.SearchReply{
		Comics: []*searchpb.Comics{
			{
				Id:  1,
				Url: "test",
			},
		},
	}, nil
}

func (m *MockClient) SearchIndex(ctx context.Context, in *searchpb.SearchRequest, opts ...grpc.CallOption) (*searchpb.SearchReply, error) {
	return &searchpb.SearchReply{
		Comics: []*searchpb.Comics{
			{
				Id:  1,
				Url: "test",
			},
		},
	}, nil
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

func TestSearchOK(t *testing.T) {

	c := Client{
		log:    DummyLogger{},
		client: &MockClient{},
	}
	_, err := c.Search(context.Background(), "test", 1)
	require.NoError(t, err)
}

func TestSearchIndexOK(t *testing.T) {
	c := Client{
		log:    DummyLogger{},
		client: &MockClient{},
	}
	_, err := c.SearchIndex(context.Background(), "test", 1)
	require.NoError(t, err)
}
