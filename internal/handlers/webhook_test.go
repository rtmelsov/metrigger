package handlers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestWebhook(t *testing.T) {

	type valueType struct {
		t      string
		name   string
		number float64
	}
	var tests = []struct {
		name       string
		expectCode int
		value      valueType
	}{{
		name:       "1",
		expectCode: 200,
		value: valueType{
			t:      "counter",
			name:   "fdsafd",
			number: 3242,
		},
	},
		{
			name:       "2",
			expectCode: 200,
			value: valueType{
				t:      "gauge",
				name:   "fdsafd",
				number: 3242,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", fmt.Sprintf("/update/%v/%v/%v", test.value.t, test.value.name, test.value.number), nil)
			w := httptest.NewRecorder()
			Webhook(w, r)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.expectCode, res.StatusCode)

		})
	}

}
