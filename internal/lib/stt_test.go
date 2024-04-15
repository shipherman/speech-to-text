package lib

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReqSTT(t *testing.T) {
	responseBody := "response"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, responseBody, http.StatusOK)
	}))
	defer server.Close()
	ConfigureSTT(server.URL, time.Second*1)
	_, err := ReqSTT([]byte("audio"))
	assert.NoError(t, err)
}
