package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Parse flags for testing
	config.ServerParseFlag()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestGetPingWebhook(t *testing.T) {
	var tests = []struct {
		name       string
		method     string
		expectCode int
	}{
		{
			name:       "ping",
			method:     "GET",
			expectCode: 200,
		},
	}
	if storage.ServerFlags.DataBaseDsn != "" {
		ts := httptest.NewServer(Webhook())
		for _, test := range tests {
			url := "/ping"
			resp := getReq(t, ts, test.method, url, nil, false)
			defer resp.Body.Close()

			require.Equal(t, test.expectCode, resp.StatusCode, fmt.Sprintf("url is %v, we want code like %v, but we got %v\r\n", url, test.expectCode, resp.StatusCode))
		}
	}

}

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

var jsonTests = []JSONTest{
	{
		name:       "1",
		action:     "update",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"delta":3242, "id":"jsonTest", "type":"counter"}`,
		value: JSONReqType{
			t:     "counter",
			name:  "jsonTest",
			delta: 3242,
		},
	},
	{
		name:       "2",
		action:     "update",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"delta":6484, "id":"jsonTest", "type":"counter"}`,
		value: JSONReqType{
			t:     "counter",
			name:  "jsonTest",
			delta: 3242,
		},
	},
	{
		name:       "3",
		action:     "update",
		expectBody: `{"id":"jsonTest", "type":"gauge", "value":32.42}`,
		method:     "POST",
		expectCode: 200,
		value: JSONReqType{
			t:     "gauge",
			name:  "jsonTest",
			value: 32.42,
		},
	},
	{
		name:       "4",
		action:     "value",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"delta":6484, "id":"jsonTest", "type":"counter"}`,
		value: JSONReqType{
			t:    "counter",
			name: "jsonTest",
		},
	},
	{
		name:       "5",
		action:     "value",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"id":"jsonTest", "type":"gauge", "value":32.42}`,
		value: JSONReqType{
			t:    "gauge",
			name: "jsonTest",
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

var gzipTests = []JSONTest{
	{
		name:       "1",
		action:     "update",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"delta":3242, "id":"gzipTest", "type":"counter"}`,
		value: JSONReqType{
			t:     "counter",
			name:  "gzipTest",
			delta: 3242,
		},
	},
	{
		name:       "2",
		action:     "update",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"delta":6484, "id":"gzipTest", "type":"counter"}`,
		value: JSONReqType{
			t:     "counter",
			name:  "gzipTest",
			delta: 3242,
		},
	},
	{
		name:       "3",
		action:     "update",
		expectBody: `{"id":"gzipTest", "type":"gauge", "value":32.42}`,
		method:     "POST",
		expectCode: 200,
		value: JSONReqType{
			t:     "gauge",
			name:  "gzipTest",
			value: 32.42,
		},
	},
	{
		name:       "4",
		action:     "value",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"delta":6484, "id":"gzipTest", "type":"counter"}`,
		value: JSONReqType{
			t:    "counter",
			name: "gzipTest",
		},
	},
	{
		name:       "5",
		action:     "value",
		method:     "POST",
		expectCode: 200,
		expectBody: `{"id":"gzipTest", "type":"gauge", "value":32.42}`,
		value: JSONReqType{
			t:    "gauge",
			name: "gzipTest",
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

func jsonReqCheck(t *testing.T, ts *httptest.Server, test *JSONTest, b *models.Metrics, isGzip bool) {
	url := fmt.Sprintf("/%s/", test.action)
	resp := getReq(t, ts, test.method, url, b, isGzip)

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	require.Equal(t, test.expectCode, resp.StatusCode, fmt.Sprintf("url is %v, we want code like %v, but we got %v, body: %v\r\n", url, test.expectCode, resp.StatusCode, string(bodyBytes)))
	// Read the response body
	require.NoError(t, err, "Error reading response body")

	var responseBody string
	if isGzip {
		var data *bytes.Buffer
		data, err = helpers.DecompressData(bodyBytes)
		assert.NoError(t, err)
		responseBody = data.String()
	} else {
		responseBody = string(bodyBytes)
	}
	// Convert response body to string

	// Use require.JSONEq to compare JSON strings
	if resp.StatusCode == http.StatusOK {
		require.JSONEq(t, test.expectBody, responseBody, "Response body does not match expected JSON")
	}
}

func TestGzipUpdateWebhook(t *testing.T) {
	ts := httptest.NewServer(Webhook())
	var b models.Metrics

	for _, test := range gzipTests {

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
		jsonReqCheck(t, ts, &test, &b, true)
	}
}

func TestJsonUpdateWebhook(t *testing.T) {
	ts := httptest.NewServer(Webhook())
	var b models.Metrics

	for _, test := range jsonTests {

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
		jsonReqCheck(t, ts, &test, &b, false)
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
		resp := getReq(t, ts, test.method, url, nil, false)
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
		resp := getReq(t, ts, test.method, test.url, nil, false)
		defer resp.Body.Close()

		require.Equal(t, test.expectCode, resp.StatusCode, fmt.Sprintf("url is %v, we want code like %v, but we got %v\r\n", test.url, test.expectCode, resp.StatusCode))
		require.Equal(t, test.contentType, resp.Header.Get("Content-Type"))
	}
}

func getReq(t *testing.T, r *httptest.Server, method, path string, body *models.Metrics, isGzip bool) *http.Response {
	var reqBody io.Reader
	if isGzip {
		data, err := json.Marshal(body)
		assert.NoError(t, err)

		res, err := helpers.CompressData(data)
		assert.NoError(t, err)
		reqBody = res
	} else if body != nil {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(body)
		assert.NoError(t, err)
		reqBody = &buf
	}

	url := r.URL + path
	req, err := http.NewRequest(method, url, reqBody)
	require.NoError(t, err)
	if isGzip {
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
	}

	resp, err := r.Client().Do(req)
	require.NoError(t, err)
	return resp
}
