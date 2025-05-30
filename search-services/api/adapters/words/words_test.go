package words

import (
	"context"
	"errors"
	"testing"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wordspb "yadro.com/course/proto/words"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in)
	return nil, args.Error(1)
}

func (m *MockClient) Norm(ctx context.Context, in *wordspb.WordsRequest, opts ...grpc.CallOption) (*wordspb.WordsReply, error) {
	if in.Phrase == "ok" {
		return &wordspb.WordsReply{
			Words: []string{"foo", "bar"},
		}, nil
	} else if in.Phrase == "error" {
		return nil, errors.New("Could not norm phrase")
	}
	return nil, nil
}

type DummyLogger struct{}

func (d DummyLogger) Error(msg string, keysAndValues ...interface{}) {fmt.Print(msg)}
func (d DummyLogger) Info(msg string, keysAndValues ...interface{})  {fmt.Print(msg)}

func TestPing(t *testing.T) {

	type expectedBehaviour func(client *MockClient, out error)

	testTable := []struct {
		name              string
		expectedBehaviour expectedBehaviour
		expectedOut       error
	}{
		{
			name: "Fail",
			expectedBehaviour: func(client *MockClient, err error) {
				client.On("Ping", mock.Anything, mock.AnythingOfType("*emptypb.Empty")).Return(nil, err)
			},
			expectedOut: errors.New("error"),
		},
		{
			name: "OK",
			expectedBehaviour: func(client *MockClient, err error) {
				client.On("Ping", mock.Anything, mock.AnythingOfType("*emptypb.Empty")).Return(nil, err)
			},
			expectedOut: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockClient := &MockClient{}
			c := Client{
				log:    DummyLogger{},
				client: mockClient,
			}
			testCase.expectedBehaviour(mockClient, testCase.expectedOut)
			err := c.Ping(context.Background())
			if testCase.expectedOut != nil {
				require.ErrorContains(t, err, testCase.expectedOut.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNorm(t *testing.T) {
	c := Client{
		log:    DummyLogger{},
		client: &MockClient{},
	}
	testTable := []struct {
		name              string
		expectedBehaviour string
		expectedOut       any
	}{
		{
			name:              "OK",
			expectedBehaviour: "ok",
			expectedOut:       []string{"foo", "bar"},
		},
		{
			name:              "Fail",
			expectedBehaviour: "error",
			expectedOut:       []string(nil),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			out, _ := c.Norm(context.Background(), testCase.expectedBehaviour)
			require.Equal(t, testCase.expectedOut, out)
		})
	}
}
