package main

import (
	"MangaReader/internal/server"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

func main() {
	log.SetDefault(log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		Level:           log.DebugLevel,
	}))

	log.Infof("Starting server on port %s", os.Getenv("PORT"))
	newServer := server.NewServer()

	err := newServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
