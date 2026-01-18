package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = r.Body.Close()

		log.Printf("Webhook received: %s\n", string(body))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	log.Println("Webhook mock running on :9090 (POST /webhook)")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}
