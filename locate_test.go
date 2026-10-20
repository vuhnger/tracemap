package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func fakeIPMap(t *testing.T, answers map[string]string, calls *int) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*calls++
		ip := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), "/best")
		switch body, ok := answers[ip]; {
		case ip == "192.168.1.254":
			http.Error(w, `{"code":"BadRequest"}`, http.StatusBadRequest)
		case ok:
			w.Write([]byte(body))
		default:
			http.Error(w, "{}", http.StatusServiceUnavailable)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func TestLocateLooksUpEachUnknownIPOnce(t *testing.T) {
	calls := 0
	server := fakeIPMap(t, map[string]string{
		"84.208.0.1":    `{"location":null}`,
		"62.115.137.62": `{"location":{"cityName":"Los Angeles","countryCodeAlpha2":"US","latitude":34.05,"longitude":-118.24,"iataCode":"LAX"}}`,
	}, &calls)
	oslo := &Location{City: "Oslo", Country: "NO", Latitude: 59.9, Longitude: 10.7}

	got, err := locate(server.URL,
		[]string{"1.1.1.1", "192.168.1.254", "84.208.0.1", "62.115.137.62", "62.115.137.62"},
		Locations{"1.1.1.1": oslo},
	)
	if err != nil {
		t.Fatal(err)
	}

	if calls != 3 {
		t.Errorf("got %d IPmap calls, want 3", calls)
	}
	want := Locations{
		"1.1.1.1":       oslo,
		"192.168.1.254": nil,
		"84.208.0.1":    nil,
		"62.115.137.62": {City: "Los Angeles", Country: "US", Latitude: 34.05, Longitude: -118.24},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestLocateFailsInsteadOfCachingAMiss(t *testing.T) {
	calls := 0
	server := fakeIPMap(t, nil, &calls)

	_, err := locate(server.URL, []string{"62.115.137.62"}, Locations{})
	if err == nil || !strings.Contains(err.Error(), "503") {
		t.Errorf("got error %v, want one naming 503", err)
	}
}
