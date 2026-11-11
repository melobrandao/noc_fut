package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFilters(t *testing.T) {
	h, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		query string
		count int
	}{{"?team=vasco", 3}, {"?team=VASCO&status=live", 1}, {"?team=Palmeiras&status=finished", 1}, {"?team=inexistente", 0}} {
		t.Run(tc.query, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest("GET", "/api/matches"+tc.query, nil))
			var body struct {
				Matches []Match `json:"matches"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if w.Code != 200 || len(body.Matches) != tc.count {
				t.Fatalf("status=%d matches=%d want=%d", w.Code, len(body.Matches), tc.count)
			}
		})
	}
}

func TestRoutes(t *testing.T) {
	h, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		method, path string
		code         int
	}{{"GET", "/", 200}, {"GET", "/api/health", 200}, {"GET", "/api/matches/1", 200}, {"GET", "/api/matches/999", 404}, {"GET", "/api/matches?status=bad", 400}, {"POST", "/api/matches", 405}, {"GET", "/api/unknown", 404}} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		if w.Code != tc.code {
			t.Errorf("%s %s: got %d want %d", tc.method, tc.path, w.Code, tc.code)
		}
	}
}

func TestLineups(t *testing.T) {
	h, err := New()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/matches/1", nil))
	var m Match
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatal(err)
	}
	for _, side := range []string{"home", "away"} {
		if len(m.Lineups[side]) != 11 {
			t.Fatalf("%s: expected 11 players", side)
		}
	}
}
