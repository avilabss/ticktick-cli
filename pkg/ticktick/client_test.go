package ticktick

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewTicktickClient_ValidToken(t *testing.T) {
	client, err := NewTicktickClient("test-token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client.APIToken != "test-token" {
		t.Errorf("expected token 'test-token', got %q", client.APIToken)
	}
	if client.BaseURL != "https://api.ticktick.com/api" {
		t.Errorf("expected default base URL, got %q", client.BaseURL)
	}
	if client.Pomodoro == nil {
		t.Error("expected Pomodoro service to be initialized")
	}
}

func TestNewTicktickClient_EmptyToken(t *testing.T) {
	_, err := NewTicktickClient("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewTicktickClient_WithHTTPClient(t *testing.T) {
	mock := &mockHTTPClient{}
	client, err := NewTicktickClient("token", WithHTTPClient(mock))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client.HTTPClient != mock {
		t.Error("expected custom HTTP client to be set")
	}
}

func TestNewTicktickClient_WithNilHTTPClient(t *testing.T) {
	_, err := NewTicktickClient("token", WithHTTPClient(nil))
	if err == nil {
		t.Fatal("expected error for nil HTTP client")
	}
}

func TestNewTicktickClient_OptionError(t *testing.T) {
	badOption := func(c *Client) error {
		return fmt.Errorf("option error")
	}
	_, err := NewTicktickClient("token", badOption)
	if err == nil {
		t.Fatal("expected error from bad option")
	}
}

func TestGet_Success(t *testing.T) {
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Cookie") != "t=test-token" {
				t.Errorf("expected cookie 't=test-token', got %q", req.Header.Get("Cookie"))
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("ok")),
			}, nil
		},
	}

	client, _ := NewTicktickClient("test-token", WithHTTPClient(mock))
	res, err := client.Get("/v2/test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}
}

func TestGet_NonOKStatus(t *testing.T) {
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Status:     "401 Unauthorized",
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	_, err := client.Get("/v2/test")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestGet_RequestError(t *testing.T) {
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("connection refused")
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	_, err := client.Get("/v2/test")
	if err == nil {
		t.Fatal("expected error for failed request")
	}
}

func TestGet_CookieHeader(t *testing.T) {
	var capturedReq *http.Request
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedReq = req
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	client, _ := NewTicktickClient("my-secret-token", WithHTTPClient(mock))
	res, _ := client.Get("/v2/test")
	defer func() { _ = res.Body.Close() }()

	cookie := capturedReq.Header.Get("Cookie")
	if cookie != "t=my-secret-token" {
		t.Errorf("expected 't=my-secret-token', got %q", cookie)
	}
}

func TestGet_URLConstruction(t *testing.T) {
	var capturedURL string
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedURL = req.URL.String()
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	res, _ := client.Get("/v2/pomodoros/timeline")
	defer func() { _ = res.Body.Close() }()

	expected := "https://api.ticktick.com/api/v2/pomodoros/timeline"
	if capturedURL != expected {
		t.Errorf("expected URL %q, got %q", expected, capturedURL)
	}
}
