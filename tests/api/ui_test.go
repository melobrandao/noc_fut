//go:build integration && ui

package api_test

import (
	"automacao-api-go/internal/dashboard"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/mxschmitt/playwright-go"
)

func TestDashboardBrowser(t *testing.T) {
	h, err := dashboard.New()
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(h)
	defer server.Close()
	pw, err := playwright.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer pw.Stop()
	browser, err := pw.Chromium.Launch()
	if err != nil {
		t.Fatal(err)
	}
	defer browser.Close()
	page, err := browser.NewPage(playwright.BrowserNewPageOptions{Viewport: &playwright.Size{Width: 1440, Height: 1100}})
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var errors []string
	page.OnPageError(func(err error) { mu.Lock(); defer mu.Unlock(); errors = append(errors, err.Error()) })
	if _, err = page.Goto(server.URL); err != nil {
		t.Fatal(err)
	}
	if err = page.Locator("[data-match='1']").WaitFor(); err != nil {
		t.Fatal(err)
	}
	count := func(selector string, want int) {
		t.Helper()
		got, err := page.Locator(selector).Count()
		if err != nil || got != want {
			t.Fatalf("%s count=%d want=%d err=%v", selector, got, want, err)
		}
	}
	click := func(selector string) {
		t.Helper()
		if err := page.Locator(selector).Click(); err != nil {
			t.Fatal(err)
		}
	}
	count(".match-card", 3)
	count(".player", 11)
	dir := filepath.Join("..", "..", "reports")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{Path: playwright.String(filepath.Join(dir, "dashboard-desktop.png")), FullPage: playwright.Bool(true)}); err != nil {
		t.Fatal(err)
	}
	click("[data-status='live']")
	count(".match-card", 1)
	click("[data-tab='stats']")
	count(".stats-row", 5)
	click("[data-tab='events']")
	count(".event", 4)
	click("[data-tab='lineup']")
	click("[data-side='away']")
	count(".player", 11)
	click("[data-favorite='1']")
	click("#nav-favorites")
	count(".match-card", 1)
	if _, err = page.Reload(); err != nil {
		t.Fatal(err)
	}
	if err = page.Locator(".star.on").WaitFor(); err != nil {
		t.Fatal(err)
	}
	values := []string{""}
	if _, err = page.Locator("#team-filter").SelectOption(playwright.SelectOptionValues{Values: &values}); err != nil {
		t.Fatal(err)
	}
	count(".match-card", 8)
	if err = page.Locator("#search").Fill("inexistente"); err != nil {
		t.Fatal(err)
	}
	count(".match-card", 0)
	count("#clear-filters", 1)
	click("#clear-filters")
	count(".match-card", 8)
	values = []string{"vasco"}
	if _, err = page.Locator("#team-filter").SelectOption(playwright.SelectOptionValues{Values: &values}); err != nil {
		t.Fatal(err)
	}
	count(".match-card", 3)
	if err = page.SetViewportSize(390, 844); err != nil {
		t.Fatal(err)
	}
	overflow, err := page.Evaluate("() => document.documentElement.scrollWidth > innerWidth")
	if err != nil || overflow == true {
		t.Fatalf("Mobile overflow=%v err=%v", overflow, err)
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{Path: playwright.String(filepath.Join(dir, "dashboard-mobile.png")), FullPage: playwright.Bool(true)}); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(errors) > 0 {
		t.Fatalf("JavaScript errors: %v", errors)
	}
}
