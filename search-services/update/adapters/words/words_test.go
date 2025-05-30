package words

import (
	"context"
	"os"
	"reflect"
	"testing"

	"log/slog"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wordspb "yadro.com/course/proto/words"
)

func TestClient_Norm(t *testing.T) {
	tests := []struct {
		name       string
		phrase     string
		wantErr    bool
		wantWords  []string
		mockClient func(*MockWordsClient)
	}{
		{
			name:      "success",
			phrase:    "hello world",
			wantErr:   false,
			wantWords: []string{"hello", "world"},
			mockClient: func(mc *MockWordsClient) {
				mc.On("Norm", context.Background(), &wordspb.WordsRequest{Phrase: "hello world"}).Return(&wordspb.WordsReply{Words: []string{"hello", "world"}}, nil)
			},
		},
		{
			name:      "error",
			phrase:    "hello world",
			wantErr:   true,
			wantWords: nil,
			mockClient: func(mc *MockWordsClient) {
				mc.On("Norm", context.Background(), &wordspb.WordsRequest{Phrase: "hello world"}).Return(&wordspb.WordsReply{Words: []string{}}, status.Errorf(codes.Internal, "internal error"))
			},
		},
		{
			name:      "empty phrase",
			phrase:    "",
			wantErr:   true,
			wantWords: nil,
			mockClient: func(mc *MockWordsClient) {
				mc.On("Norm", context.Background(), &wordspb.WordsRequest{Phrase: ""}).Return(&wordspb.WordsReply{Words: []string{}}, status.Errorf(codes.Internal, "internal error"))
			},
		},
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})
	logger := slog.New(handler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &MockWordsClient{}
			tt.mockClient(mc)

			c := &Client{
				client: mc,
				log:    logger,
			}

			got, err := c.Norm(context.Background(), tt.phrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Norm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantWords) {
				t.Errorf("Client.Norm() = %v, want %v", got, tt.wantWords)
			}
		})
	}
}

type MockWordsClient struct {
	mock.Mock
}

func (m *MockWordsClient) Norm(ctx context.Context, in *wordspb.WordsRequest, opts ...grpc.CallOption) (*wordspb.WordsReply, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*wordspb.WordsReply), args.Error(1)
}

func (m *MockWordsClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	if in == nil {
		return nil, nil
	}
	args := m.Called(ctx, in)
	return nil, args.Error(1)
}

func TestService_Ping(t *testing.T) {
	mc := &MockWordsClient{}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})
	logger := slog.New(handler)

	c := &Client{
		client: mc,
		log:    logger,
	}

	mc.On("Ping", context.Background(), nil).Return(emptypb.Empty{}, nil)

	t.Run("returns no error", func(t *testing.T) {
		err := c.Ping(context.Background())
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}
