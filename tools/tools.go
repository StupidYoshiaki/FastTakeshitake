package tools

import (
	"log"
	"net/http"
	"os"
	"time"
)

func StartPingRoutine() {
	interval := 1 * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	port := os.Getenv("PORT")
	url := "http://localhost:" + port + "/"

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		<-ticker.C
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("Ping failed: %v", err)
			continue
		}
		resp.Body.Close()
		log.Printf("Pinged %s successfully", url)
	}
}
