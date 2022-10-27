package models

import "time"

type Empty struct{}

type Info struct {
	Interval time.Duration    `json:"interval"`
	URLs     map[string]Empty `json:"urls"`
}

type Result struct {
	Date time.Time       `json:"date"`
	URLs map[string]bool `json:"urls"`
}

type Results struct {
	Results []Result `json:"results"`
}
