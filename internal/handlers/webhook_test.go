package handlers

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostWebhook(t *testing.T) {

	type valueType struct {
		t      string
		name   string
		number float64
	}
	var tests = []struct {
		name       string
		method     string
		expectCode int
		value      valueType
	}{{
		name:       "1",
		method:     "POST",
		expectCode: 200,
		value: valueType{
			t:      "counter",
			name:   "fdsafd",
			number: 3242,
		},
	},
		{
			name:       "2",
			method:     "POST",
			expectCode: 200,
			value: valueType{
				t:      "gauge",
				name:   "fdsafd",
				number: 3242,
			},
		},
	}
	ts := httptest.NewServer(Webhook())
	for _, test := range tests {
		url := fmt.Sprintf("/update/%v/%v/%v", test.value.t, test.value.name, test.value.number)
		resp := getReq(t, ts, test.method, url)

		require.Equal(t, test.expectCode, resp.StatusCode, fmt.Sprintf("url is %v, we want code like %v, but we got %v\r\n", url, test.expectCode, resp.StatusCode))
	}
}

func TestGetWebhook(t *testing.T) {

	var tests = []struct {
		name        string
		method      string
		contentType string
		expectCode  int
		url         string
	}{
		{
			name:        "1",
			method:      "GET",
			contentType: "text/plain; charset=utf-8",
			expectCode:  200,
			url:         "/value/counter/fdsafd",
		},
		{
			name:        "2",
			method:      "GET",
			contentType: "text/plain; charset=utf-8",
			expectCode:  200,
			url:         "/value/gauge/fdsafd",
		},
		{
			name:        "3",
			method:      "GET",
			contentType: "text/html",
			expectCode:  200,
			url:         "/",
		},
	}
	ts := httptest.NewServer(Webhook())
	for _, test := range tests {
		resp := getReq(t, ts, test.method, test.url)

		require.Equal(t, test.expectCode, resp.StatusCode, fmt.Sprintf("url is %v, we want code like %v, but we got %v\r\n", test.url, test.expectCode, resp.StatusCode))
		require.Equal(t, test.contentType, resp.Header.Get("Content-Type"))
	}
}

func getReq(t *testing.T, r *httptest.Server, method, path string) *http.Response {
	url := r.URL + path
	req, err := http.NewRequest(method, url, nil)
	require.NoError(t, err)

	resp, err := r.Client().Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp

}
