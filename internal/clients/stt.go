// Interaction with STT service
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

// temporary logger
var logEntry = middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}

func ConfigureSTT(a string, t time.Duration) {
	cfg.Address = a
	cfg.Timeout = t
}

// Request STT for transcribtion
func ReqSTT(filepath string) error {
	var wavFileBytes []byte
	client := resty.New()
	client.RetryMaxWaitTime = time.Second * 1
	client.RetryCount = 5

	// read file
	wavFileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Build connection string for STT app
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

		fmt.Printf("resp code: %v; resp body: %v; Addr: %s\n", resp.StatusCode(), resp, reqAddress)

		switch resp.StatusCode() {
		case 200:
			fmt.Println(string(resp.Body()))
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
	err = backoff.Retry(f, b)
	if err != nil {
		logEntry.Logger.Print(fmt.Errorf("ReqSTT error: %w", err))
		return err
	}

	return nil

}
