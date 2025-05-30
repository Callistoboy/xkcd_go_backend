package core

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type Service struct {
	log         *slog.Logger
	db          DB
	xkcd        XKCD
	words       Words
	concurrency int
	inProgress  atomic.Bool
	lock        sync.Mutex
}

func NewService(
	log *slog.Logger, db DB, xkcd XKCD, words Words, concurrency int,
) (*Service, error) {
	if concurrency < 1 {
		return nil, fmt.Errorf("wrong concurrency specified: %d", concurrency)
	}
	return &Service{
		log:         log,
		db:          db,
		xkcd:        xkcd,
		words:       words,
		concurrency: concurrency,
	}, nil
}

func (s *Service) Update(ctx context.Context) (err error) {
	if ok := s.lock.TryLock(); !ok {
		s.log.Error("Service already runs update")
		return ErrAlreadyExists
	}
	defer s.lock.Unlock()

	s.inProgress.Store(true)
	defer s.inProgress.Store(false)

	s.log.Info("Update started")
	defer func(start time.Time) {
		s.log.Info("Update finished", "duration", time.Since(start), "error", err)
	}(time.Now())

	IDs, err := s.db.IDs(ctx)
	if err != nil {
		s.log.Error("Failed to get existing IDs in DB", "error", err)
		return fmt.Errorf("failed to get existing IDs in DB: %v", err)
	}

	s.log.Debug("Existing comics in DB", "count", len(IDs))
	exists := make(map[int]bool, len(IDs))
	for _, id := range IDs {
		exists[id] = true
	}

	// get last comics ID
	latestComicID, err := s.xkcd.LastID(ctx)
	if err != nil {
		s.log.Error("Failed to get latest comic", "error", err)
		return fmt.Errorf("failed to get last ID in XKCD: %v", err)
	}

	s.log.Debug("Last comics ID in XKCD", "id", latestComicID)

	generator := generateIDs(ctx, 1, latestComicID, exists)
	fetchers := s.getComics(ctx, generator)

	var errorsFound bool
	var added int
	for info := range fetchers {
		words, err := s.words.Norm(ctx, info.Description+" "+info.Title)
		if err != nil {
			errorsFound = true
			s.log.Error("Failed to normalize", "id", info.ID, "error", err)
			continue
		}
		err = s.db.Add(ctx, Comics{
			ID:    info.ID,
			URL:   info.URL,
			Words: words,
		})
		if err != nil {
			errorsFound = true
			s.log.Error("Failed to save comics", "id", info.ID, "error", err)
			continue
		}
		added++
	}
	s.log.Debug("Added new comics", "count", added)

	if errorsFound {
		return fmt.Errorf("failed to fetch/store some comics")
	}

	return nil
}

func generateIDs(ctx context.Context, first, last int, exists map[int]bool) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := first; i <= last; i++ {
			if exists[i] {
				continue
			}
			select {
			case <-ctx.Done():
				return
			default:
				ch <- i
			}
		}
	}()
	return ch
}

func (s *Service) getComics(ctx context.Context, in <-chan int) <-chan XKCDInfo {
	out := make(chan XKCDInfo)
	var wg sync.WaitGroup
	wg.Add(s.concurrency)

	for i := range s.concurrency {
		go func() {
			s.log.Debug("Fetcher up", "id", i)
			defer s.log.Debug("Fetcher down", "id", i)
			defer wg.Done()
			for id := range in {
				if id == 404 {
					// special case
					out <- XKCDInfo{ID: id, Title: "404", Description: "Not found"}
					continue
				}
				info, err := s.xkcd.Get(ctx, id)
				if err != nil {
					s.log.Error("Failed to get comics", "id", id, "error", err)
					continue
				}
				s.log.Debug("Fetched", "id", id)
				out <- info
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (s *Service) Stats(ctx context.Context) (ServiceStats, error) {
	dbStats, err := s.db.Stats(ctx)
	if err != nil {
		s.log.Error("failed to get stats", "error", err)
		return ServiceStats{}, err
	}
	lastID, err := s.xkcd.LastID(ctx)
	if err != nil {
		s.log.Error("failed to get last comics ID", "error", err)
		return ServiceStats{}, err
	}
	return ServiceStats{
		DBStats:     dbStats,
		ComicsTotal: lastID,
	}, nil
}

func (s *Service) Status(ctx context.Context) ServiceStatus {
	if s.inProgress.Load() {
		return StatusRunning
	}
	return StatusIdle

}

func (s *Service) Drop(ctx context.Context) error {
	err := s.db.Drop(ctx)
	if err != nil {
		s.log.Error("failed to drop db entries", "error", err)
	}
	s.log.Debug("DB dropped")
	return err
}
