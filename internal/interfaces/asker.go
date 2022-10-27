package interfaces

import (
	"context"
	"github.com/MikhailSolovev/URLAsker/internal/models"
	"time"
)

// This interface can be use as public contract

type Asker interface {
	GetInfo(ctx context.Context) (info models.Info, err error)
	// ListResults - list results between two dates
	ListResults(ctx context.Context, dateFrom, dateTo time.Time) (result models.Results, err error)
	SetInterval(ctx context.Context, interval time.Duration) (err error)
	// SetURLs - rewrite urls in set
	SetURLs(ctx context.Context, urls ...string) (err error)
	// AddURLs - append urls to set
	AddURLs(ctx context.Context, urls ...string) (err error)
	// DeleteURLs - delete urls from set
	DeleteURLs(ctx context.Context, urls ...string) (err error)
}
