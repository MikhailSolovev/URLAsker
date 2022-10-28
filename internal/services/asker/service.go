package asker

import (
	"context"
	"github.com/MikhailSolovev/URLAsker/internal/models"
	"net/http"
	"time"
)

// Service TODO: replace info to some persistence config storage (file, Redis)
type Service struct {
	db   Storage
	info *models.Info
}

type Storage interface {
	ListLatestResult(ctx context.Context) (results models.ResultsPostgresDTO, err error)
	ListResults(ctx context.Context, dateFrom, dateTo time.Time) (results models.ResultsPostgresDTO, err error)
	RecordResult(ctx context.Context, result models.ResultPostgresDTO) (err error)
}

func New(db Storage, info *models.Info) *Service {
	return &Service{db: db, info: info}
}

func (s *Service) Run(ctx context.Context) (err error) {
	var data models.ResultPostgresDTO
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.info.Ticker.C:
			data.Date = time.Now().UTC()
			for url := range s.info.URLs {
				data.URL = url
				res, err := http.Get(url)
				if err != nil {
					data.Available = false
				} else {
					if res.StatusCode/100 == 2 {
						data.Available = true
					} else {
						data.Available = false
					}
				}
				if err = s.db.RecordResult(context.Background(), data); err != nil {
					return err
				}
			}
		}
	}
}

func (s *Service) GetInfo(ctx context.Context) (info models.Info, err error) {
	s.info.RLock()
	defer s.info.RUnlock()
	return *s.info, nil
}

func (s *Service) ListLatestResult(ctx context.Context) (result models.Result, err error) {
	data, err := s.db.ListLatestResult(ctx)
	if err != nil {
		return
	}

	result = models.Result{URLs: map[string]bool{}}
	for i, dtoResult := range data.Results {
		if i == 0 {
			result.Date = dtoResult.Date
		}
		result.URLs[dtoResult.URL] = dtoResult.Available
	}

	return
}

func (s *Service) ListResults(ctx context.Context, dateFrom, dateTo time.Time) (results models.Results, err error) {
	data, err := s.db.ListResults(ctx, dateFrom, dateTo)
	if err != nil {
		return
	}

	var result models.Result
	for i, dtoResult := range data.Results {
		if result.Date != dtoResult.Date {
			if i != 0 {
				results.Results = append(results.Results, result)
			}
			result = models.Result{Date: dtoResult.Date, URLs: map[string]bool{}}
			result.URLs[dtoResult.URL] = dtoResult.Available
		} else {
			result.URLs[dtoResult.URL] = dtoResult.Available
		}
	}
	return
}

func (s *Service) SetInterval(ctx context.Context, interval time.Duration) (err error) {
	s.info.Lock()
	defer s.info.Unlock()
	s.info.Interval = interval
	s.info.Ticker.Reset(s.info.Interval)
	return
}

func (s *Service) SetURLs(ctx context.Context, urls ...string) (err error) {
	s.info.Lock()
	defer s.info.Unlock()
	s.info.URLs = map[string]models.Empty{}
	for _, url := range urls {
		s.info.URLs[url] = models.Empty{}
	}
	return
}

func (s *Service) AddURLs(ctx context.Context, urls ...string) (err error) {
	s.info.Lock()
	defer s.info.Unlock()
	for _, url := range urls {
		if _, ok := s.info.URLs[url]; !ok {
			s.info.URLs[url] = models.Empty{}
		}
	}
	return
}

func (s *Service) DeleteURLs(ctx context.Context, urls ...string) (err error) {
	s.info.Lock()
	defer s.info.Unlock()
	for _, url := range urls {
		delete(s.info.URLs, url)
	}
	return
}
