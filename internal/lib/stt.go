// Interaction with STT service [docker image with STT model deployed]
package clients

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-resty/resty/v2"
)

type options struct {
	Address string
	Timeout time.Duration
}

// STT config
var cfg = options{}

// Transcibed text
var text string

// temporary logger
var logEntry = middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}

func ConfigureSTT(a string, t time.Duration) {
	cfg.Address = a
	cfg.Timeout = t
}

// Request STT for transcribtion
func ReqSTT(wavFileBytes []byte) (string, error) {
	// Init htttp client to connect to STT server
	// To Do
	// Move to init()
	client := resty.New()
	client.RetryMaxWaitTime = time.Second * 1
	client.RetryCount = 5

	// Build connection string for STT app
	// TO DO
	// Move to init()
	reqAddress := fmt.Sprintf("%s/stt", cfg.Address)

	// Create lambda to use it in backoff.Retry()
	f := func() error {
		// Post wav file to STT  service
		resp, err := client.R().EnableTrace().
			SetHeader("Content-Type", "application/octet-stream").
			SetBody(wavFileBytes).
			Post(reqAddress)
		if err != nil {
			return err
		}

		// fmt.Printf("resp code: %v; resp body: %v; Addr: %s\n", resp.StatusCode(), resp, reqAddress)

		switch resp.StatusCode() {
		case 200:
			// fmt.Println(string(resp.Body()))
			text = resp.String()
		// Server error
		case 500:
			fmt.Println(string(resp.Body()))
		default:
			return nil
		}

		return nil
	}

	// Use backoff package to implement retryer with increasing interval between attempts
	b := backoff.NewExponentialBackOff()
	b.MaxInterval = cfg.Timeout
	err := backoff.Retry(f, b)
	if err != nil {
		logEntry.Logger.Print(fmt.Errorf("ReqSTT error: %w", err))
		return text, err
	}

	return text, nil

}
