package handlers_test

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"

	"github.com/rtmelsov/metrigger/internal/handlers"
)

// ExamplePingDBHandler демонстрирует, как использовать PingDBHandler.
//
// Этот пример создаёт тестовый HTTP-запрос и записывает ответ в буфер.
func ExamplePingDBHandler() {
	// создаём новый рекордер для записи ответа
	w := httptest.NewRecorder()

	// создаём тестовый запрос (метод GET, без тела)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	// вызываем сам обработчик
	handlers.PingDBHandler(w, req)

	// печатаем код ответа
	fmt.Println(w.Code)

	// печатаем тело ответа
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// ok
}

// ExampleMetricsListHandler демонстрирует, как использовать MetricsListHandler.
//
// Этот пример создаёт тестовый HTTP-запрос и записывает ответ в буфер.
func ExampleMetricsListHandler() {
	// создаём новый рекордер для записи ответа
	w := httptest.NewRecorder()

	// создаём тестовый запрос (метод GET, без тела)
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// вызываем сам обработчик
	handlers.MetricsListHandler(w, req)

	// печатаем код ответа
	fmt.Println(w.Code)

	// Output:
	// 200
}

// ExampleJSONUpdate демонстрирует, как использовать JSONUpdate.
func ExampleJSONUpdate() {
	// создаём новый рекордер для записи ответа
	w := httptest.NewRecorder()

	// создаём тестовый запрос (метод POST, тело - типа models.Metrics)
	body := []byte(`{
    "type": "gauge",
    "id":    "gauge",
    "value": 3212
}`)
	req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// вызываем сам обработчик
	handlers.JSONUpdate(w, req)

	// печатаем код ответа
	fmt.Println(w.Code)

	// печатаем тело ответа
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// object
}

// ExampleJSONGet демонстрирует, как использовать JSONGet.
func ExampleJSONGet() {
	// создаём новый рекордер для записи ответа
	w := httptest.NewRecorder()

	// создаём тестовый запрос (метод POST, тело - типа models.Metrics)
	body := []byte(`{
		"type": "gauge",
		"id":    "gauge"
	}`)
	req := httptest.NewRequest(http.MethodPost, "/value", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// вызываем сам обработчик
	handlers.JSONGet(w, req)

	// печатаем код ответа
	fmt.Println(w.Code)

	// печатаем тело ответа
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// object
}

// ExampleMetricsUpdateHandler в примере создаётся роутер, регистрируются маршруты и отправляется тестовый запрос.
func ExampleMetricsUpdateHandler() {
	// 1. Создаём новый роутер
	r := chi.NewRouter()

	// 2. Регистрируем хендлеры
	handlers.MetricsValueHandler(r)

	// 3. Создаём тестовый запрос (например, GET /gauge/myMetric)
	req := httptest.NewRequest(http.MethodGet, "/update/gauge/myMetric/232", nil)
	w := httptest.NewRecorder()

	// 4. Пускаем запрос через роутер
	r.ServeHTTP(w, req)

	// 5. Печатаем код ответа
	fmt.Println(w.Code)

	// 6. Печатаем тело ответа
	fmt.Println(w.Body.String())

	// Output:
	// 404
	// can't find parameters
}

// ExampleMetricsValueHandler в примере создаётся роутер, регистрируются маршруты и отправляется тестовый запрос.
func ExampleMetricsValueHandler() {
	// 1. Создаём новый роутер
	r := chi.NewRouter()

	// 2. Регистрируем хендлеры
	handlers.MetricsValueHandler(r)

	// 3. Создаём тестовый запрос (например, GET /gauge/myMetric)
	req := httptest.NewRequest(http.MethodGet, "/value/gauge/myMetric", nil)
	w := httptest.NewRecorder()

	// 4. Пускаем запрос через роутер
	r.ServeHTTP(w, req)

	// 5. Печатаем код ответа
	fmt.Println(w.Code)

	// 6. Печатаем тело ответа
	fmt.Println(w.Body.String())

	// Output:
	// 200
}
