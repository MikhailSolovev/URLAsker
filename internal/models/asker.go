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

// swagger:response infoResponse
type InfoRestDTO struct {
	// example: 10s
	Interval string `json:"interval"`
	// example: ["https://www.google.com/", "https://yandex.com/"]
	URLs []string `json:"urls"`
}

// swagger:response listLatestResponse
type Result struct {
	// example: 2022-10-30T20:04:11.827677Z
	Date time.Time `json:"date"`
	// example: {"http://xydsffew.com/": false, "https://www.google.com/": true}
	URLs map[string]bool `json:"urls"`
}

// swagger:response listResponse
type Results struct {
	// example: {
	//  "results": [
	//    {
	//      "date": "2022-10-30T20:04:06.840497Z",
	//      "urls": {
	//        "http://xydsffew.com/": false,
	//        "https://www.google.com/": true,
	//        "https://www.yahoo.com/": true,
	//        "https://yandex.com/": true
	//      }
	//    }
	//    }
	//  ]
	//}
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
