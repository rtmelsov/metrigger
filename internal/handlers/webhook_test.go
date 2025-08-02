package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/constants"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Parse flags for testing
	config.ServerParseFlag()

	filePath := storage.ServerFlags.FileStoragePath // Change this to the file you want to delete

	err := os.Remove(filePath)
	if err != nil {
		fmt.Println("Error deleting file:", err)
		//return
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestGetPingWebhook(t *testing.T) {
	if storage.ServerFlags.DataBaseDsn != "" {
		r, err := Webhook()
		require.NoError(t, err, "Error to get routers")
		ts := httptest.NewServer(r)
		for _, test := range constants.GetPingWebhook {
			url := "/ping"
			resp := getReq(t, ts, test.Method, url, nil, false)
			defer resp.Body.Close()

			require.Equal(t, test.ExpectCode, resp.StatusCode, fmt.Sprintf("url is %v, we want code like %v, but we got %v\r\n", url, test.ExpectCode, resp.StatusCode))
		}
	}
}

func jsonReqCheck(t *testing.T, ts *httptest.Server, test *models.JSONTest, b *models.Metrics, isGzip bool) {
	url := fmt.Sprintf("/%s/", test.Action)
	resp := getReq(t, ts, test.Method, url, b, isGzip)

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	require.Equal(t, test.ExpectCode, resp.StatusCode, fmt.Sprintf("url is %v, we want code like %v, but we got %v, body: %v\r\n", url, test.ExpectCode, resp.StatusCode, string(bodyBytes)))
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

	// Use require.JSONEq to compare JSON strings
	if resp.StatusCode == http.StatusOK {
		fmt.Println("expected body", test.ExpectBody)
		fmt.Println("response body", responseBody)

		require.JSONEq(t, test.ExpectBody, responseBody, "Response body does not match expected JSON")
	}
}

func TestGzipUpdateWebhook(t *testing.T) {
	r, err := Webhook()
	require.NoError(t, err, "Error to get routers")
	ts := httptest.NewServer(r)
	var b models.Metrics

	for _, test := range constants.GzipTests {

		if test.Value.T == "counter" {
			b = models.Metrics{
				MType: test.Value.T,
				ID:    test.Value.Name,
				Delta: &test.Value.Delta,
			}
		} else {
			b = models.Metrics{
				MType: test.Value.T,
				ID:    test.Value.Name,
				Value: &test.Value.Value,
			}
		}
		jsonReqCheck(t, ts, &test, &b, true)
	}
}

func TestJSONUpdateWebhook(t *testing.T) {
	r, err := Webhook()
	require.NoError(t, err, "Error to get routers")

	ts := httptest.NewServer(r)
	var b models.Metrics

	for _, test := range constants.JSONTests {

		if test.Value.T == "counter" {
			b = models.Metrics{
				MType: test.Value.T,
				ID:    test.Value.Name,
				Delta: &test.Value.Delta,
			}
		} else {
			b = models.Metrics{
				MType: test.Value.T,
				ID:    test.Value.Name,
				Value: &test.Value.Value,
			}
		}
		jsonReqCheck(t, ts, &test, &b, false)
	}
}

func TestPostWebhook(t *testing.T) {
	r, err := Webhook()
	require.NoError(t, err, "Error to get routers")

	ts := httptest.NewServer(r)

	for _, test := range constants.PostWebhook {
		url := fmt.Sprintf("/update/%v/%v/%v", test.Value.T, test.Value.Name, test.Value.Number)
		resp := getReq(t, ts, test.Method, url, nil, false)
		defer resp.Body.Close()

		require.Equal(t, test.ExpectCode, resp.StatusCode, fmt.Sprintf("url is %v, we want code like %v, but we got %v\r\n", url, test.ExpectCode, resp.StatusCode))
	}
}

func TestGetWebhook(t *testing.T) {
	r, err := Webhook()
	require.NoError(t, err, "Error to get routers")

	ts := httptest.NewServer(r)
	for _, test := range constants.GetWebhook {
		resp := getReq(t, ts, test.Method, test.URL, nil, false)
		defer resp.Body.Close()

		require.Equal(t, test.ExpectCode, resp.StatusCode, fmt.Sprintf("url is %v, we want code like %v, but we got %v\r\n", test.URL, test.ExpectCode, resp.StatusCode))
		require.Equal(t, test.ContentType, resp.Header.Get("Content-Type"))
	}
}

func TestGetHtmlWebhook(t *testing.T) {
	r, err := Webhook()
	require.NoError(t, err, "Error to get routers")

	ts := httptest.NewServer(r)

	resp := getReq(t, ts, "GET", "", nil, false)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode, fmt.Sprintf("we want code like %v, but we got %v\r\n", 200, resp.StatusCode))

	require.Equal(t, "text/html", resp.Header.Get("Content-Type"))

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

func jsonBenchReqCheck(ts *httptest.Server, test *models.JSONTest, b *models.Metrics, isGzip bool) error {
	url := fmt.Sprintf("/%s/", test.Action)
	return getBenchReq(ts, test.Method, url, b, isGzip)
}

func BenchmarkGzipUpdateWebhook(ben *testing.B) {
	r, err := Webhook()
	require.NoError(ben, err, "Error to get routers")

	ts := httptest.NewServer(r)
	var b models.Metrics

	for i := 0; i < ben.N; i++ {
		for _, test := range constants.GzipTests {

			if test.Value.T == "counter" {
				b = models.Metrics{
					MType: test.Value.T,
					ID:    test.Value.Name,
					Delta: &test.Value.Delta,
				}
			} else {
				b = models.Metrics{
					MType: test.Value.T,
					ID:    test.Value.Name,
					Value: &test.Value.Value,
				}
			}
			err := jsonBenchReqCheck(ts, &test, &b, true)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func BenchmarkJSONUpdateWebhook(ben *testing.B) {
	r, err := Webhook()
	require.NoError(ben, err, "Error to get routers")

	ts := httptest.NewServer(r)
	var b models.Metrics

	for i := 0; i < ben.N; i++ {
		for _, test := range constants.JSONTests {

			if test.Value.T == "counter" {
				b = models.Metrics{
					MType: test.Value.T,
					ID:    test.Value.Name,
					Delta: &test.Value.Delta,
				}
			} else {
				b = models.Metrics{
					MType: test.Value.T,
					ID:    test.Value.Name,
					Value: &test.Value.Value,
				}
			}
			err := jsonBenchReqCheck(ts, &test, &b, false)
			if err != nil {
				log.Println(err)

			}
		}
	}
}

func BenchmarkPostWebhook(ben *testing.B) {
	r, err := Webhook()
	require.NoError(ben, err, "Error to get routers")

	ts := httptest.NewServer(r)
	for i := 0; i < ben.N; i++ {
		for _, test := range constants.PostWebhook {
			url := fmt.Sprintf("/update/%v/%v/%v", test.Value.T, test.Value.Name, test.Value.Number)
			err := getBenchReq(ts, test.Method, url, nil, false)
			if err != nil {
				log.Println(err)
			}

		}
	}
}

func BenchmarkGetWebhook(ben *testing.B) {
	r, err := Webhook()
	require.NoError(ben, err, "Error to get routers")

	ts := httptest.NewServer(r)
	for i := 0; i < ben.N; i++ {
		for _, test := range constants.GetWebhook {
			err := getBenchReq(ts, test.Method, test.URL, nil, false)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func BenchmarkGetHtmlWebhook(ben *testing.B) {
	r, err := Webhook()
	require.NoError(ben, err, "Error to get routers")

	ts := httptest.NewServer(r)
	for i := 0; i < ben.N; i++ {
		err := getBenchReq(ts, "GET", "", nil, false)
		if err != nil {
			log.Println(err)
		}
	}
}

func getBenchReq(r *httptest.Server, method, path string, body *models.Metrics, isGzip bool) error {
	var reqBody io.Reader
	if isGzip {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}

		res, err := helpers.CompressData(data)
		if err != nil {
			return err
		}
		reqBody = res
	} else if body != nil {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return err
		}
		reqBody = &buf
	}

	url := r.URL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return err
	}
	if isGzip {
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
	}

	resp, err := r.Client().Do(req)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
