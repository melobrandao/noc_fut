//go:build integration

package api_test

import (
	"automacao-api-go/internal/dashboard"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/mxschmitt/playwright-go"
)

// TestAPIWithPlaywright exercises HTTP requests through APIRequestContext.
func TestAPIWithPlaywright(t *testing.T) {
	h, err := dashboard.New()
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(h)
	defer server.Close()
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("Playwright: %v. Execute scripts/setup-playwright.ps1", err)
	}
	defer pw.Stop()
	request, err := pw.Request.NewContext(playwright.APIRequestNewContextOptions{BaseURL: playwright.String(server.URL)})
	if err != nil {
		t.Fatal(err)
	}
	defer request.Dispose()
	for _, tc := range []struct {
		path   string
		status int
		count  int
	}{
		{"/api/matches?team=vasco", 200, 3},
		{"/api/matches?team=vasco&status=live", 200, 1},
		{"/api/matches?team=inexistente", 200, 0},
		{"/api/matches?status=invalid", 400, -1},
		{"/api/matches/999", 404, -1},
	} {
		t.Run(tc.path, func(t *testing.T) {
			response, err := request.Get(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Dispose()
			if response.Status() != tc.status {
				t.Fatalf("status=%d want=%d", response.Status(), tc.status)
			}
			if tc.count < 0 {
				return
			}
			body, err := response.Body()
			if err != nil {
				t.Fatal(err)
			}
			var payload struct {
				Matches []dashboard.Match `json:"matches"`
			}
			if err := json.Unmarshal(body, &payload); err != nil {
				t.Fatal(err)
			}
			if len(payload.Matches) != tc.count {
				t.Fatalf("matches=%d want=%d", len(payload.Matches), tc.count)
			}
			for _, m := range payload.Matches {
				if m.Home.ID != "vasco" && m.Away.ID != "vasco" {
					t.Fatal("Filtro retornou partida sem Vasco")
				}
			}
		})
	}
}
