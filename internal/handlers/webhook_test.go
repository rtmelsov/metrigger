package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type JSONReqType struct {
	t     string
	name  string
	delta int64
	value float64
}

type JSONTest struct {
	name       string
	method     string
	expectBody string
	expectCode int
	action     string
	value      JSONReqType
}

var updateTests = []JSONTest{
	{
		name:       "1",
		action:     "update",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"Type":"counter","Value":3242}`,
		value: JSONReqType{
			t:     "counter",
			name:  "fdsafd",
			delta: 3242,
		},
	},
	{
		name:       "2",
		action:     "update",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"Type":"counter","Value":6484}`,
		value: JSONReqType{
			t:     "counter",
			name:  "fdsafd",
			delta: 3242,
		},
	},
	{
		name:       "3",
		action:     "update",
		expectBody: `{"Type":"gauge","Value":32.42}`,
		method:     "POST",
		expectCode: 200,
		value: JSONReqType{
			t:     "gauge",
			name:  "fdsafd",
			value: 32.42,
		},
	},
	{
		name:       "4",
		action:     "value",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"Type":"counter","Value":6484}`,
		value: JSONReqType{
			t:    "counter",
			name: "fdsafd",
		},
	},
	{
		name:       "5",
		action:     "value",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"Type":"gauge","Value":32.42}`,
		value: JSONReqType{
			t:    "gauge",
			name: "fdsafd",
		},
	},
	{
		name:       "6",
		action:     "value",
		method:     "POST",
		expectCode: 404,
		expectBody: "",
		value: JSONReqType{
			t:    "gauge",
			name: "unknown",
		},
	},
}

func jsonReqCheck(t *testing.T, ts *httptest.Server, test *JSONTest, b *models.Metrics) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(b)
	assert.NoError(t, err)
	url := fmt.Sprintf("/%s/", test.action)
	resp := getReq(t, ts, test.method, url, &buf)

	defer resp.Body.Close()

	require.Equal(t, test.expectCode, resp.StatusCode, fmt.Sprintf("url is %v, we want code like %v, but we got %v\r\n", url, test.expectCode, resp.StatusCode))
	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Error reading response body")

	// Convert response body to string
	responseBody := string(bodyBytes)

	fmt.Println("responseBody", responseBody)

	// Use require.JSONEq to compare JSON strings
	if resp.StatusCode == http.StatusOK {
		require.JSONEq(t, test.expectBody, responseBody, "Response body does not match expected JSON")
	}
}

func TestJsonUpdateWebhook(t *testing.T) {
	ts := httptest.NewServer(Webhook())
	var b models.Metrics

	for _, test := range updateTests {

		if test.value.t == "counter" {
			b = models.Metrics{
				MType: test.value.t,
				ID:    test.value.name,
				Delta: &test.value.delta,
			}
		} else {
			b = models.Metrics{
				MType: test.value.t,
				ID:    test.value.name,
				Value: &test.value.value,
			}
		}
		jsonReqCheck(t, ts, &test, &b)
	}
}

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
		resp := getReq(t, ts, test.method, url, nil)
		defer resp.Body.Close()

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
		resp := getReq(t, ts, test.method, test.url, nil)
		defer resp.Body.Close()

		require.Equal(t, test.expectCode, resp.StatusCode, fmt.Sprintf("url is %v, we want code like %v, but we got %v\r\n", test.url, test.expectCode, resp.StatusCode))
		require.Equal(t, test.contentType, resp.Header.Get("Content-Type"))
	}
}

func getReq(t *testing.T, r *httptest.Server, method, path string, body io.Reader) *http.Response {
	var reqBody io.Reader
	if body != nil {
		reqBody = body
	}
	url := r.URL + path
	req, err := http.NewRequest(method, url, reqBody)
	require.NoError(t, err)

	resp, err := r.Client().Do(req)
	require.NoError(t, err)
	return resp
}
