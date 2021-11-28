package websiteracer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T) {
	t.Run("return fastest url", func (t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)
		slowUrl := slowServer.URL
		fastUrl := fastServer.URL

		want := fastUrl
		got, _ := Racer(slowUrl, fastUrl) 

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}

		slowServer.Close()
		fastServer.Close()
	})

	t.Run("returns an error if the server doesn't respond in 10s", func(t *testing.T) {
		server := makeDelayedServer(25 * time.Millisecond)
		defer server.Close()

		_, err := ConfigurableRacer(server.URL, server.URL, 20 * time.Millisecond)

		if err == nil {
			t.Error("expected an error but didn't get one")

		}
	})
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}