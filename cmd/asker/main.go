package main

import (
	"fmt"
	"github.com/MikhailSolovev/URLAsker/internal/config"
	"github.com/MikhailSolovev/URLAsker/pkg/logger"
	"log"
)

func main() {
	cfg := config.New()
	if cfg == nil {
		log.Fatal("Can't parse config!")
	}

	lg := logger.New(cfg.ServiceName, cfg.IsPretty, cfg.LogLevel)
	lg.LogInfo(fmt.Sprintf("Logger init for %v with %v level", cfg.ServiceName, cfg.LogLevel))
}
