package mocks

import (
	"PriemBot/scraper/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type MockBrowserlessClient struct {
	ExecuteScriptFunc func(script string) ([]models.ScraperResultItem, error)
}

func (m *MockBrowserlessClient) ExecuteScript(script string) ([]models.ScraperResultItem, error) {
	return m.ExecuteScriptFunc(script)
}

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func NewMockHTTPResponse(statusCode int, body interface{}) *http.Response {
	jsonBody, _ := json.Marshal(body)
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBuffer(jsonBody)),
	}
}
