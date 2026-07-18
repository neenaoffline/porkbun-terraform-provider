package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateDNSRecordParsesStringID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dns/create/example.com" {
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"SUCCESS","id":"123456789"}`)
	}))
	defer server.Close()

	client := NewClient("api-key", "secret-key")
	client.baseURL = server.URL

	id, err := client.CreateDNSRecord("example.com", CreateDNSRecordRequest{
		Type:    "A",
		Content: "192.0.2.1",
	})
	if err != nil {
		t.Fatalf("CreateDNSRecord returned an error: %v", err)
	}
	if id != "123456789" {
		t.Fatalf("unexpected record ID: got %q", id)
	}
}

func TestCreateDNSRecordReportsAPIErrorBeforeParsingObjectID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ERROR","message":"record already exists","id":{"detail":"duplicate"}}`)
	}))
	defer server.Close()

	client := NewClient("api-key", "secret-key")
	client.baseURL = server.URL

	_, err := client.CreateDNSRecord("example.com", CreateDNSRecordRequest{
		Type:    "A",
		Content: "192.0.2.1",
	})
	if err == nil {
		t.Fatal("CreateDNSRecord returned no error")
	}
	if !strings.Contains(err.Error(), "record already exists") {
		t.Fatalf("unexpected error: %v", err)
	}
}
