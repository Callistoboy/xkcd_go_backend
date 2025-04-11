package core

import (
	"cmp"
	"context"
	"log/slog"
	"maps"
	"slices"
)

type Service struct {
	log   *slog.Logger
	db    DB
	words Words
}

func NewService(
	log *slog.Logger, db DB, words Words) (*Service, error) {
	return &Service{
		log:   log,
		db:    db,
		words: words,
	}, nil
}

func (s *Service) Ping(ctx context.Context) error {
	return nil
}

func (s *Service) Search(ctx context.Context, phrase string, limit int) ([]Comics, error) {
	normed, err := s.words.Norm(ctx, phrase)
	if err != nil {
		s.log.Error("failed to find keywords", "error", err)
		return nil, err
	}
	s.log.Debug("normalized query", "keywords", normed)

	scores := map[int]int{}
	for _, keyword := range normed {
		IDs, err := s.db.Search(ctx, keyword)
		if err != nil {
			s.log.Error("failed to search keyword in DB", "error", err)
			return nil, err
		}
		for _, ID := range IDs {
			scores[ID]++
		}
	}
	s.log.Debug("relevant comics", "count", len(scores))
	sorted := slices.SortedFunc(maps.Keys(scores), func(a, b int) int {
		return cmp.Compare(scores[b], scores[a])
	})

	if len(sorted) < limit {
		limit = len(sorted)
	}
	sorted = sorted[:limit]

	result := make([]Comics, 0, len(sorted))
	for _, id := range sorted {
		comic, err := s.db.Get(ctx, id)
		if err != nil {
			s.log.Error("failed to get comic", "id", id, "error", err)
			return nil, err
		}
		result = append(result, comic)
	}

	s.log.Debug("returning comics", "count", len(result))

	return result, nil
}
