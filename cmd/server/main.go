package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jonasqasoftware/quickmetrics-project/internal/httpapi"
	"github.com/jonasqasoftware/quickmetrics-project/web"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           httpapi.New(web.Static()),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("quickmetrics listening on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
