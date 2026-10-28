package vaultkey

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWalletCreatePreservesProviderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/sdk/wallets" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Retry-After", "3600")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"wallet creation rate limit exceeded"}`))
	}))
	defer server.Close()
	client, err := NewClient("testnet_key", "testnet_secret", WithBaseURL(server.URL+"/sdk"))
	if err != nil {
		t.Fatal(err)
	}
	_, providerErr, err := client.Wallets.Create(context.Background(), CreateWalletPayload{UserID: "user", ChainType: ChainTypeEVM})
	if err != nil {
		t.Fatal(err)
	}
	if providerErr == nil || providerErr.StatusCode != 429 || providerErr.RetryAfter != time.Hour || !strings.Contains(providerErr.Message, "rate limit") {
		t.Fatalf("provider error = %#v", providerErr)
	}
}
