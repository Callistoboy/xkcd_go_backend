package core

import (
	"context"
	"errors"
	"os"

	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWords struct {
	mock.Mock
}

func (m *MockWords) Norm(ctx context.Context, phrase string) ([]string, error) {
	args := m.Called(ctx, phrase)
	return args.Get(0).([]string), args.Error(1)
}

func TestService_Ping(t *testing.T) {
	s := &Service{}

	t.Run("returns no error", func(t *testing.T) {
		err := s.Ping(context.Background())
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("returns no error with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := s.Ping(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func TestSearchIndex(t *testing.T) {
	tests := []struct {
		name        string
		phrase      string
		limit       int
		wordsErr    error
		indexResult []int
		wantErr     bool
	}{
		{
			name:        "valid phrase and limit",
			phrase:      "test phrase",
			limit:       10,
			wordsErr:    nil,
			indexResult: []int{1, 2, 3},
			wantErr:     false,
		},
		{
			name:        "invalid phrase",
			phrase:      "invalid phrase",
			limit:       10,
			wordsErr:    errors.New("norm error"),
			indexResult: nil,
			wantErr:     true,
		},
		{
			name:        "empty phrase",
			phrase:      "",
			limit:       10,
			wordsErr:    nil,
			indexResult: nil,
			wantErr:     false,
		},
		{
			name:        "zero limit",
			phrase:      "test phrase",
			limit:       0,
			wordsErr:    nil,
			indexResult: []int{1, 2, 3},
			wantErr:     false,
		},
		{
			name:        "nil context",
			phrase:      "test phrase",
			limit:       10,
			wordsErr:    nil,
			indexResult: []int{1, 2, 3},
			wantErr:     false,
		},
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})
	logger := slog.New(handler)

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			mockWords := &MockWords{}

			mockWords.On("Norm", mock.Anything, tt.phrase).Return([]string{"keyword1", "keyword2"}, tt.wordsErr)

			s := &Service{
				log:   logger,
				db:    nil,
				words: mockWords,
				index: NewIndex(),
			}

			ctx := context.Background()
			if tt.name == "nil context" {
				ctx = nil
			}

			_, err := s.SearchIndex(ctx, tt.phrase, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Search(ctx context.Context, keyword string) ([]int, error) {
	args := m.Called(ctx, keyword)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockDB) Get(ctx context.Context, id int) (Comics, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Comics), args.Error(1)
}

func (m *MockDB) LastID(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return 1, args.Error(1)
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name       string
		phrase     string
		limit      int
		wordsErr   error
		dbErr      error
		wantErr    bool
		wantResult []Comics
	}{
		{
			name:       "successful search",
			phrase:     "test phrase",
			limit:      10,
			wordsErr:   nil,
			dbErr:      nil,
			wantErr:    false,
			wantResult: []Comics{{ID: 1, URL: "url1", Score: 1}, {ID: 2, URL: "url2", Score: 1}},
		},
		{
			name:       "invalid phrase",
			phrase:     "invalid phrase",
			limit:      10,
			wordsErr:   errors.New("invalid phrase"),
			dbErr:      nil,
			wantErr:    true,
			wantResult: nil,
		},
		{
			name:       "invalid keyword",
			phrase:     "test phrase",
			limit:      10,
			wordsErr:   nil,
			dbErr:      errors.New("invalid keyword"),
			wantErr:    true,
			wantResult: nil,
		},
		{
			name:       "empty phrase",
			phrase:     "",
			limit:      10,
			wordsErr:   nil,
			dbErr:      nil,
			wantErr:    false,
			wantResult: make([]Comics, 0),
		},
		{
			name:       "limit 0",
			phrase:     "test phrase",
			limit:      0,
			wordsErr:   nil,
			dbErr:      nil,
			wantErr:    false,
			wantResult: make([]Comics, 0),
		},
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})
	logger := slog.New(handler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := &MockWords{}
			db := &MockDB{}
			s := &Service{words: words, db: db, log: logger}

			if tt.phrase != "" {
				words.On("Norm", mock.Anything, tt.phrase).Return([]string{"keyword1", "keyword2"}, tt.wordsErr)
			} else {
				words.On("Norm", mock.Anything, tt.phrase).Return([]string{}, tt.wordsErr)
			}

			db.On("Get", mock.Anything, 1).Return(Comics{ID: 1, URL: "url1"}, tt.dbErr)
			db.On("Get", mock.Anything, 2).Return(Comics{ID: 2, URL: "url2"}, tt.dbErr)
			db.On("Get", mock.Anything, 3).Return(Comics{ID: 3, URL: "url3"}, tt.dbErr)
			db.On("Get", mock.Anything, 4).Return(Comics{ID: 4, URL: "url4"}, tt.dbErr)
			db.On("Search", mock.Anything, "keyword1").Return([]int{1}, tt.dbErr)
			db.On("Search", mock.Anything, "keyword2").Return([]int{2}, tt.dbErr)

			got, err := s.Search(context.Background(), tt.phrase, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.wantResult, got) {
				t.Errorf("Service.Search() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}
