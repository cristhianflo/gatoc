package subcommands

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchRatesArrayResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"moneda":"USD","fuente":"oficial","nombre":"Dólar","promedio":737.8816,"fechaActualizacion":"2026-07-23T00:00:00-04:00"}]`))
	}))
	defer server.Close()

	result := fetchRates(&http.Client{Timeout: 2 * time.Second}, server.URL)
	if result.Err != nil {
		t.Fatalf("expected no error, got %v", result.Err)
	}

	if len(result.Rates) != 1 {
		t.Fatalf("expected 1 rate, got %d", len(result.Rates))
	}

	if result.Rates[0].Currency != "USD" {
		t.Fatalf("expected USD currency, got %q", result.Rates[0].Currency)
	}
}

func TestFetchRatesSingleObjectResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"moneda":"EUR","fuente":"oficial","nombre":"Euro","promedio":841.84173862,"fechaActualizacion":"2026-07-23T00:00:00-04:00"}`))
	}))
	defer server.Close()

	result := fetchRates(&http.Client{Timeout: 2 * time.Second}, server.URL)
	if result.Err != nil {
		t.Fatalf("expected no error, got %v", result.Err)
	}

	if len(result.Rates) != 1 {
		t.Fatalf("expected 1 rate, got %d", len(result.Rates))
	}

	if result.Rates[0].Currency != "EUR" {
		t.Fatalf("expected EUR currency, got %q", result.Rates[0].Currency)
	}
}

func TestFetchRatesFailureResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":"upstream"}`))
	}))
	defer server.Close()

	result := fetchRates(&http.Client{Timeout: 2 * time.Second}, server.URL)
	if result.Err == nil {
		t.Fatal("expected error for non-200 response")
	}
}

func TestFormatRateFields(t *testing.T) {
	rate := DolarResponse{
		Currency:  "EUR",
		Source:    "oficial",
		Name:      "Euro",
		Average:   841.84173862,
		UpdatedAt: "2026-07-23T00:00:00-04:00",
	}

	fields, err := formatRateFields(rate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}

	if fields[0].Name == "" || fields[1].Name == "" {
		t.Fatal("expected non-empty field names")
	}
}
