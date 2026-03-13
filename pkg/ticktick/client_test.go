package ticktick

import (
	"encoding/json"
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

func TestPost_Success(t *testing.T) {
	var capturedReq *http.Request
	var capturedBody string
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedReq = req
			body, _ := io.ReadAll(req.Body)
			capturedBody = string(body)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	res, err := client.Post("/v2/test", strings.NewReader(`{"key":"value"}`))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	if capturedReq.Method != "POST" {
		t.Errorf("expected POST, got %s", capturedReq.Method)
	}
	if capturedReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", capturedReq.Header.Get("Content-Type"))
	}
	if capturedReq.Header.Get("Cookie") != "t=token" {
		t.Errorf("expected cookie 't=token', got %q", capturedReq.Header.Get("Cookie"))
	}
	if capturedBody != `{"key":"value"}` {
		t.Errorf("expected body '{\"key\":\"value\"}', got %q", capturedBody)
	}
}

func TestPut_Success(t *testing.T) {
	var capturedMethod string
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedMethod = req.Method
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	res, err := client.Put("/v2/test", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	if capturedMethod != "PUT" {
		t.Errorf("expected PUT, got %s", capturedMethod)
	}
}

func TestDelete_Success(t *testing.T) {
	var capturedMethod string
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedMethod = req.Method
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	res, err := client.Delete("/v2/test", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	if capturedMethod != "DELETE" {
		t.Errorf("expected DELETE, got %s", capturedMethod)
	}
}

func TestDelete_Accepts204(t *testing.T) {
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNoContent,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	res, err := client.Delete("/v2/test", nil)
	if err != nil {
		t.Fatalf("expected no error for 204, got %v", err)
	}
	defer func() { _ = res.Body.Close() }()
}

func TestPostJSON_Success(t *testing.T) {
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"id2etag":{"abc":"xyz"},"id2error":{}}`)),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))

	reqBody := map[string]string{"key": "value"}
	var result BatchResponse
	err := client.PostJSON("/v2/test", reqBody, &result)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID2Etag["abc"] != "xyz" {
		t.Errorf("expected etag 'xyz', got %q", result.ID2Etag["abc"])
	}
}

func TestPostJSON_NilResult(t *testing.T) {
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	err := client.PostJSON("/v2/test", map[string]string{}, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPostJSON_VerifiesRequestBody(t *testing.T) {
	var capturedBody map[string]any
	mock := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			_ = json.Unmarshal(body, &capturedBody)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil
		},
	}

	client, _ := NewTicktickClient("token", WithHTTPClient(mock))
	type testReq struct {
		Title string `json:"title"`
		Count int    `json:"count"`
	}
	_ = client.PostJSON("/v2/test", testReq{Title: "hello", Count: 42}, nil)

	if capturedBody["title"] != "hello" {
		t.Errorf("expected title 'hello', got %v", capturedBody["title"])
	}
	if capturedBody["count"] != float64(42) {
		t.Errorf("expected count 42, got %v", capturedBody["count"])
	}
}
