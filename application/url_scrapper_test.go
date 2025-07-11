package application

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang_test_task1/domain"
)

type mockUnavailableHostsTransport struct {
	defaultTransport http.RoundTripper
	unavailableHosts map[string]bool
	availableHosts   map[string]bool
}

func (m *mockUnavailableHostsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	host := req.URL.Hostname()
	if m.unavailableHosts[host] {
		return nil, fmt.Errorf("host %s is unreachable", host)
	}

	if m.availableHosts[host] {
		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	return m.defaultTransport.RoundTrip(req)
}

func TestUrlScrapper_GetInfoByURL_AllCases(t *testing.T) {
	testClient := &http.Client{
		Transport: &mockUnavailableHostsTransport{
			defaultTransport: http.DefaultTransport,
			unavailableHosts: map[string]bool{
				"unavailable-domain.example": true,
			},
			availableHosts: map[string]bool{
				"external.com": true,
			},
		},
	}

	tests := []struct {
		name           string
		serverResponse string
		statusCode     int
		expectError    bool
		expectedInfo   *domain.WebsiteInfo
		serverDelay    time.Duration
	}{
		{
			name: "Successful analysis",
			serverResponse: `
				<!DOCTYPE html>
				<html>
				<head><title>Test Page</title></head>
				<body>
					<h1>Heading 1</h1>
					<h2>Heading 2</h2>
					<a href="/internal">Internal Link</a>
					<a href="/broken-link">Internal broken Link</a>
					<a href="http://external.com">External Link</a>
					<a href="http://unavailable-domain.example">External Unavailable Link</a>
					<form><input type="password" /></form>
				</body>
				</html>`,
			statusCode:  http.StatusOK,
			expectError: false,
			expectedInfo: &domain.WebsiteInfo{
				HTMLVersion:       "HTML5",
				Title:             "Test Page",
				HeadingsCounts:    map[string]int{"h1": 1, "h2": 1, "h3": 0, "h4": 0, "h5": 0, "h6": 0},
				InternalLinks:     2,
				ExternalLinks:     2,
				InaccessibleLinks: 2,
				IsExistLoginForm:  true,
			},
		},
		{
			name:           "Empty response",
			serverResponse: "",
			statusCode:     http.StatusOK,
			expectedInfo: &domain.WebsiteInfo{
				HTMLVersion:       "Unknown",
				Title:             "",
				HeadingsCounts:    map[string]int{"h1": 0, "h2": 0, "h3": 0, "h4": 0, "h5": 0, "h6": 0},
				InternalLinks:     0,
				ExternalLinks:     0,
				InaccessibleLinks: 0,
				IsExistLoginForm:  false,
			},
		},
		{
			name:           "URL not found",
			serverResponse: "",
			statusCode:     http.StatusNotFound,
			expectError:    true,
		},
		{
			name:           "Timeout error",
			serverResponse: "",
			statusCode:     http.StatusOK,
			expectError:    true,
			serverDelay:    11 * time.Second, // Simulate a timeout
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/broken-link" {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				time.Sleep(tt.serverDelay)
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			scrapper := NewUrlScrapper(testClient)

			info, err := scrapper.GetInfoByURL(context.Background(), server.URL)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if info.Title != tt.expectedInfo.Title {
					t.Errorf("expected title '%s', got '%s'", tt.expectedInfo.Title, info.Title)
				}

				if info.HTMLVersion != tt.expectedInfo.HTMLVersion {
					t.Errorf("expected HTML version '%s', got '%s'", tt.expectedInfo.HTMLVersion, info.HTMLVersion)
				}

				if info.InternalLinks != tt.expectedInfo.InternalLinks || info.ExternalLinks != tt.expectedInfo.ExternalLinks {
					t.Errorf("unexpected link counts: internal=%d, external=%d", info.InternalLinks, info.ExternalLinks)
				}

				if info.InaccessibleLinks != tt.expectedInfo.InaccessibleLinks {
					t.Errorf("unexpected inaccessible link counts: %d", info.InaccessibleLinks)
				}

				if info.IsExistLoginForm != tt.expectedInfo.IsExistLoginForm {
					t.Errorf("expected login form existence to be %v", tt.expectedInfo.IsExistLoginForm)
				}

				for heading, count := range tt.expectedInfo.HeadingsCounts {
					if info.HeadingsCounts[heading] != count {
						t.Errorf("expected %d headings of type '%s', got %d", count, heading, info.HeadingsCounts[heading])
					}
				}
			}
		})
	}
}
