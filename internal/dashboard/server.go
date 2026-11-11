// Package dashboard serves the local football dashboard and its demo API.
package dashboard

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
)

//go:embed web/* fixtures.json
var files embed.FS

type Team struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
	Color string `json:"color"`
}

type Match struct {
	ID        int                 `json:"id"`
	Home      Team                `json:"home"`
	Away      Team                `json:"away"`
	Status    string              `json:"status"`
	Time      string              `json:"time"`
	Minute    int                 `json:"minute"`
	HomeScore int                 `json:"homeScore"`
	AwayScore int                 `json:"awayScore"`
	Stadium   string              `json:"stadium"`
	Lineups   map[string][]Player `json:"lineups"`
	Stats     []Stat              `json:"stats"`
	Events    []Event             `json:"events"`
}

type Player struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

type Stat struct {
	Label string `json:"label"`
	Home  int    `json:"home"`
	Away  int    `json:"away"`
}

type Event struct {
	Minute int    `json:"minute"`
	Type   string `json:"type"`
	Team   string `json:"team"`
	Text   string `json:"text"`
}

// New returns an isolated handler. All data is fictional and bundled locally.
func New() (http.Handler, error) {
	data, err := files.ReadFile("fixtures.json")
	if err != nil {
		return nil, err
	}
	var matches []Match
	if err = json.Unmarshal(data, &matches); err != nil {
		return nil, err
	}
	assets, err := fs.Sub(files, "web")
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if !getOnly(w, r) {
			return
		}
		writeJSON(w, map[string]string{"status": "ok", "mode": "demo"})
	})
	mux.HandleFunc("/api/teams", func(w http.ResponseWriter, r *http.Request) {
		if !getOnly(w, r) {
			return
		}
		teams := []Team{}
		seen := map[string]bool{}
		for _, m := range matches {
			for _, t := range []Team{m.Home, m.Away} {
				if !seen[t.ID] {
					teams = append(teams, t)
					seen[t.ID] = true
				}
			}
		}
		writeJSON(w, teams)
	})
	mux.HandleFunc("/api/matches", func(w http.ResponseWriter, r *http.Request) {
		if !getOnly(w, r) {
			return
		}
		team := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("team")))
		status := r.URL.Query().Get("status")
		if status != "" && status != "live" && status != "scheduled" && status != "finished" {
			http.Error(w, "Status inválido", http.StatusBadRequest)
			return
		}
		result := []Match{}
		for _, m := range matches {
			if team != "" && !strings.Contains(strings.ToLower(m.Home.ID+" "+m.Home.Name+" "+m.Away.ID+" "+m.Away.Name), team) {
				continue
			}
			if status != "" && m.Status != status {
				continue
			}
			result = append(result, m)
		}
		writeJSON(w, map[string]any{"mode": "demo", "matches": result})
	})
	mux.HandleFunc("/api/matches/", func(w http.ResponseWriter, r *http.Request) {
		if !getOnly(w, r) {
			return
		}
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/matches/"))
		if err == nil {
			for _, m := range matches {
				if m.ID == id {
					writeJSON(w, m)
					return
				}
			}
		}
		http.NotFound(w, r)
	})
	mux.HandleFunc("/api/", http.NotFound)
	mux.Handle("/", http.FileServer(http.FS(assets)))
	return mux, nil
}

func getOnly(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(v)
}
