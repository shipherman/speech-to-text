package handlers

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/shipherman/speech-to-text/gen/ent"
)

// Put request to queue
// Check file
// Save file
// Transcribe
// Nofify client
func HandleSTT(w http.ResponseWriter, r *http.Request) {

	client, err := ent.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=postgres password=pass sslmode=disable")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	client.Schema.Create(r.Context())
	defer client.Close()

	// Save tuple to db
	_, err = client.Audio.Create().
		SetPath("path").
		SetHash("hash").
		SetText("text").
		Save(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save file
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	af, err := os.Create("/tmp/stt/audio.wav")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	af.Write(buf.Bytes())
	af.Close()
}
