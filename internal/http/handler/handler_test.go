package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/vandi37/Calculator/internal/http/handler"
	"github.com/vandi37/Calculator/pkg/logger"
)

func TestOk(t *testing.T) {
	testCases := []struct {
		req      string
		expected string
	}{
		{`{"expression":"1+1"}"`, `{"result":2}`},
		{`{"expression":"2*3+(2+2)"}`, `{"result":10}`},
	}
	handler := handler.NewHandler("/", logger.New(os.Stderr))
	for _, testCase := range testCases {
		t.Run(testCase.req, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(testCase.req)))
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			if resp.Code != http.StatusOK {
				t.Errorf("Wrong status code, expected 200, got: %d", resp.Code)
			} else if resp.Body.String() != testCase.expected {
				t.Errorf("wrong data, expected '%s', but got '%s'", testCase.expected, resp.Body.String())
			}
		})
	}
}

func TestCalcError(t *testing.T) {
	testCases := []string{
		`{"expression":"1+"}"`,
		`{"expression":"2*+(2+2)"}`,
		`{"expression":"1+1*(1"}`,
		`{"expression":"1 + b"}`,
	}

	handler := handler.NewHandler("/", logger.New(os.Stderr))
	for _, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(testCase)))
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			if resp.Code != 422 {
				t.Errorf("Wrong status code, expected 422, got: %d", resp.Code)
			}
		})
	}
}

func TestWrongBody(t *testing.T) {
	testCases := []string{
		`{"expression":1+1"}"`,
		``,
		`{"task":"1+1"}`,
		`{expression:"1 + 1"}`,
		`{"expression":1 + 1}`,
		`{"expression":"1 + 1",}`,
		`{"expression":}`,
		`{}`,
	}

	handler := handler.NewHandler("/", logger.New(os.Stderr))
	for _, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(testCase)))
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			if resp.Code != http.StatusBadRequest {
				t.Errorf("Wrong status code, expected 400, got: %d", resp.Code)
			}
		})
	}
}

func TestWrongPage(t *testing.T) {
	handler := handler.NewHandler("/", logger.New(os.Stderr))
	req := httptest.NewRequest(http.MethodPost, "/be", bytes.NewBuffer([]byte("")))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusNotFound {
		t.Errorf("Wrong status code, expected 404, got: %d", resp.Code)
	}
}

func TestWrongMethod(t *testing.T) {
	handler := handler.NewHandler("/", logger.New(os.Stderr))
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte("")))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusMethodNotAllowed {
		t.Errorf("Wrong status code, expected 405, got: %d", resp.Code)
	}
}