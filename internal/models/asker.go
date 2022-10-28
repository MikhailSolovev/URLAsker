package models

import (
	"sync"
	"time"
)

type Empty struct{}

type Info struct {
	sync.RWMutex
	Interval time.Duration
	Ticker   *time.Ticker
	URLs     map[string]Empty
}

type InfoRestDTO struct {
	Interval string   `json:"interval"`
	URLs     []string `json:"urls"`
}

type Result struct {
	Date time.Time       `json:"date"`
	URLs map[string]bool `json:"urls"`
}

type Results struct {
	Results []Result `json:"results"`
}

type ResultPostgresDTO struct {
	Date      time.Time
	URL       string
	Available bool
}

type ResultsPostgresDTO struct {
	Results []ResultPostgresDTO
}
