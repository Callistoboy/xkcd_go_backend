package grpc

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/types/known/emptypb"
	searchpb "yadro.com/course/proto/search"
	"yadro.com/course/search/core"
)

type MockSearcher struct {
	mock.Mock
}

func (m *MockSearcher) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockSearcher) Search(ctx context.Context, phrase string, limit int) ([]core.Comics, error) {
	args := m.Called(ctx, phrase, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]core.Comics), args.Error(1)
}

func (m *MockSearcher) SearchIndex(ctx context.Context, phrase string, limit int) ([]core.Comics, error) {
	args := m.Called(ctx, phrase, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]core.Comics), args.Error(1)
}

func (m *MockSearcher) BuildIndex(ctx context.Context) error {
	return nil
}

func TestPing(t *testing.T) {

	type expectedBehaviour func(client *MockSearcher, out error)

	testTable := []struct {
		name              string
		expectedBehaviour expectedBehaviour
		expectedOut       error
	}{
		{
			name: "Fail",
			expectedBehaviour: func(client *MockSearcher, err error) {
				client.On("Ping", context.Background()).Return(err)
			},
			expectedOut: errors.New("error"),
		},
		{
			name: "OK",
			expectedBehaviour: func(client *MockSearcher, err error) {
				client.On("Ping", context.Background()).Return(nil)
			},
			expectedOut: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockClient := &MockSearcher{}
			s := Server{
				service: mockClient,
			}
			testCase.expectedBehaviour(mockClient, testCase.expectedOut)
			_, err := s.Ping(context.Background(), &emptypb.Empty{})
			if testCase.expectedOut != nil {
				require.ErrorContains(t, err, testCase.expectedOut.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSearch(t *testing.T) {

	type expectedBehaviour func(client *MockSearcher)

	type args struct {
		request *searchpb.SearchRequest
	}

	testTable := []struct {
		name              string
		args              args
		expectedBehaviour expectedBehaviour
		expectedOut       any
		expectedError     error
	}{
		{
			name: "OK",
			args: args{
				request: &searchpb.SearchRequest{
					Phrase: "foo bar",
					Limit:  1,
				},
			},
			expectedBehaviour: func(client *MockSearcher) {
				client.On("Search", context.Background(), "foo bar", 1).Return([]core.Comics{
					{
						ID:       1,
						URL:      "http://foo.bar",
						Keywords: []string{"foo", "bar"},
						Score:    1,
					},
					{
						ID:       2,
						URL:      "http://abra.cadabra",
						Keywords: []string{"abra", "cadabra"},
						Score:    2,
					},
				}, nil)
			},
			expectedOut: &searchpb.SearchReply{Comics: []*searchpb.Comics{
				{
					Id:  1,
					Url: "http://foo.bar",
				},
				{
					Id:  2,
					Url: "http://abra.cadabra",
				},
			},
			},
		},
		{
			name: "OK (no limit)",
			args: args{
				request: &searchpb.SearchRequest{
					Phrase: "foo bar",
					Limit:  0,
				},
			},
			expectedBehaviour: func(client *MockSearcher) {
				client.On("Search", context.Background(), "foo bar", 10).Return([]core.Comics{
					{
						ID:       1,
						URL:      "http://foo.bar",
						Keywords: []string{"foo", "bar"},
						Score:    1,
					},
					{
						ID:       2,
						URL:      "http://abra.cadabra",
						Keywords: []string{"abra", "cadabra"},
						Score:    2,
					},
				}, nil)
			},
			expectedOut: &searchpb.SearchReply{Comics: []*searchpb.Comics{
				{
					Id:  1,
					Url: "http://foo.bar",
				},
				{
					Id:  2,
					Url: "http://abra.cadabra",
				},
			},
			},
		},
		{
			name: "Fail",
			args: args{
				request: &searchpb.SearchRequest{
					Phrase: "foo bar",
					Limit:  1,
				},
			},
			expectedBehaviour: func(client *MockSearcher) {
				client.On("Search", context.Background(), "foo bar", 1).Return(nil, errors.New("error"))
			},
			expectedError: errors.New("error"),
		},
		{
			name: "Not Founds",
			args: args{
				request: &searchpb.SearchRequest{
					Phrase: "foo bar",
					Limit:  1,
				},
			},
			expectedBehaviour: func(client *MockSearcher) {
				client.On("Search", context.Background(), "foo bar", 1).Return(nil, core.ErrNotFound)
			},
			expectedError: status.Error(codes.NotFound, "nothing found"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockClient := &MockSearcher{}
			s := Server{
				service: mockClient,
			}
			testCase.expectedBehaviour(mockClient)
			res, err := s.Search(context.Background(), testCase.args.request)
			if err != nil {
				require.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				require.Equal(t, res, testCase.expectedOut)
			}

		})
	}
}

func TestSearchIndex(t *testing.T) {

	type expectedBehaviour func(client *MockSearcher)

	type args struct {
		request *searchpb.SearchRequest
	}

	testTable := []struct {
		name              string
		args              args
		expectedBehaviour expectedBehaviour
		expectedOut       any
		expectedError     error
	}{
		{
			name: "OK",
			args: args{
				request: &searchpb.SearchRequest{
					Phrase: "foo bar",
					Limit:  1,
				},
			},
			expectedBehaviour: func(client *MockSearcher) {
				client.On("SearchIndex", context.Background(), "foo bar", 1).Return([]core.Comics{
					{
						ID:       1,
						URL:      "http://foo.bar",
						Keywords: []string{"foo", "bar"},
						Score:    1,
					},
					{
						ID:       2,
						URL:      "http://abra.cadabra",
						Keywords: []string{"abra", "cadabra"},
						Score:    2,
					},
				}, nil)
			},
			expectedOut: &searchpb.SearchReply{Comics: []*searchpb.Comics{
				{
					Id:  1,
					Url: "http://foo.bar",
				},
				{
					Id:  2,
					Url: "http://abra.cadabra",
				},
			},
			},
		},
		{
			name: "OK (no limit)",
			args: args{
				request: &searchpb.SearchRequest{
					Phrase: "foo bar",
					Limit:  0,
				},
			},
			expectedBehaviour: func(client *MockSearcher) {
				client.On("SearchIndex", context.Background(), "foo bar", 10).Return([]core.Comics{
					{
						ID:       1,
						URL:      "http://foo.bar",
						Keywords: []string{"foo", "bar"},
						Score:    1,
					},
					{
						ID:       2,
						URL:      "http://abra.cadabra",
						Keywords: []string{"abra", "cadabra"},
						Score:    2,
					},
				}, nil)
			},
			expectedOut: &searchpb.SearchReply{Comics: []*searchpb.Comics{
				{
					Id:  1,
					Url: "http://foo.bar",
				},
				{
					Id:  2,
					Url: "http://abra.cadabra",
				},
			},
			},
		},
		{
			name: "Fail",
			args: args{
				request: &searchpb.SearchRequest{
					Phrase: "foo bar",
					Limit:  1,
				},
			},
			expectedBehaviour: func(client *MockSearcher) {
				client.On("SearchIndex", context.Background(), "foo bar", 1).Return(nil, errors.New("error"))
			},
			expectedError: errors.New("error"),
		},
		{
			name: "Not Founds",
			args: args{
				request: &searchpb.SearchRequest{
					Phrase: "foo bar",
					Limit:  1,
				},
			},
			expectedBehaviour: func(client *MockSearcher) {
				client.On("SearchIndex", context.Background(), "foo bar", 1).Return(nil, core.ErrNotFound)
			},
			expectedError: status.Error(codes.NotFound, "nothing found"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockClient := &MockSearcher{}
			s := Server{
				service: mockClient,
			}
			testCase.expectedBehaviour(mockClient)
			res, err := s.SearchIndex(context.Background(), testCase.args.request)
			if err != nil {
				require.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				require.Equal(t, res, testCase.expectedOut)
			}

		})
	}
}
