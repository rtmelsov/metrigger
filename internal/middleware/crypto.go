package middleware

import (
	"bytes"
	"crypto/rsa"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/storage"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func CryptoParser(privateKey *rsa.PrivateKey) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cr := storage.ServerFlags.CryptoRate
			contentEncoding := r.Header.Get("X-Encrypted")
			bodyEncrypted := strings.Contains(contentEncoding, "true")
			if cr != "" && bodyEncrypted {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				err = r.Body.Close() // release resources
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				bodyBytes, err = helpers.DecryptFromClient(privateKey, bodyBytes)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				// 3) Replace r.Body with the modified bytes
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

				// (Optional) update Content-Length so downstream sees correct length
				r.ContentLength = int64(len(bodyBytes))
				r.Header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))

				// If you need GetBody (for retries, etc.), you can also set:
				r.GetBody = func() (io.ReadCloser, error) {
					return io.NopCloser(bytes.NewReader(bodyBytes)), nil
				}
			}

			// передаём управление хендлеры
			h.ServeHTTP(w, r)
		})
	}
}
