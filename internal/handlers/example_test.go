package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rtmelsov/metrigger/internal/agent"
	"github.com/rtmelsov/metrigger/internal/metrics"
	"github.com/rtmelsov/metrigger/internal/models"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/rtmelsov/metrigger/internal/handlers"
)

var met *models.MetricsCollector

// ExampleJSONUpdate демонстрирует, как использовать JSONUpdate.
//
// Отправка данных в виде json объектов отдельно
func ExampleJSONUpdate() {
	var pollCount float64

	// Получаем метрики runtime в виде списка
	met, _ = metrics.CollectMetric(&pollCount)

	// Перебираем значения для отдельной отправки
	for k, b := range *met {
		// создаём новый рекордер для записи ответа
		w := httptest.NewRecorder()

		// Есть два типа метрики
		//
		// Counter новое значение должно добавляться к предыдущему
		// если какое-то значение уже было известно серверу
		counter := agent.RequestToServer("counter", k, 0, 1)

		res, err := json.Marshal(*counter)
		if err != nil {
			log.Println("Error", err.Error())
		}
		req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewBuffer(res))
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

		//
		// Gauge новое значение должно замещать предыдущее
		gauge := agent.RequestToServer("gauge", k, b, 0)

		res, err = json.Marshal(*gauge)
		if err != nil {
			log.Println("Error", err.Error())
		}
		req = httptest.NewRequest(http.MethodPost, "/update", bytes.NewBuffer(res))
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
}

// ExampleJSONGet демонстрирует, как использовать JSONGet.
func ExampleJSONGet() {

	for k := range *met {
		// создаём новый рекордер для записи ответа
		w := httptest.NewRecorder()

		var counterMetric *models.Metrics
		var gaugeMetric *models.Metrics

		counterMetric = &models.Metrics{
			MType: "counter",
			ID:    k,
		}
		gaugeMetric = &models.Metrics{
			MType: "gauge",
			ID:    k,
		}

		// Создаем объект для поиска метрики в сервере counter
		res, err := json.Marshal(counterMetric)
		if err != nil {
			log.Println("Error", err.Error())
			return
		}
		req := httptest.NewRequest(http.MethodPost, "/value", bytes.NewBuffer(res))
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

		// Создаем объект для поиска метрики в сервере gauge
		res, err = json.Marshal(gaugeMetric)
		if err != nil {
			log.Println("Error", err.Error())
			return
		}
		req = httptest.NewRequest(http.MethodPost, "/value", bytes.NewBuffer(res))
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

}

// ExampleMetricsUpdateHandler в примере создаётся роутер, регистрируются маршруты и отправляется тестовый запрос.
func ExampleMetricsUpdateHandler() {

	for k, b := range *met {
		// 1. Создаём новый роутер
		r := chi.NewRouter()

		// 2. Регистрируем хендлеры
		handlers.MetricsValueHandler(r)

		// 3. Создаём тестовый запрос (например, GET /gauge/myMetric)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update/gauge/%s/%f", k, b), nil)
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

		// 3. Создаём тестовый запрос (например, GET /counter/myMetric)
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update/gauge/%s/%b", k, 1), nil)
		w = httptest.NewRecorder()

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
}

// ExampleMetricsValueHandler в примере создаётся роутер, регистрируются маршруты и отправляется тестовый запрос.
func ExampleMetricsValueHandler() {
	for k := range *met {
		// 1. Создаём новый роутер
		r := chi.NewRouter()

		// 2. Регистрируем хендлеры
		handlers.MetricsValueHandler(r)

		// 3. Создаём тестовый запрос (например, GET /gauge/myMetric)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update/gauge/%s", k), nil)
		w := httptest.NewRecorder()

		// 4. Пускаем запрос через роутер
		r.ServeHTTP(w, req)

		// 5. Печатаем код ответа
		fmt.Println(w.Code)

		// 6. Печатаем тело ответа
		fmt.Println(w.Body.String())

		// Output:
		// 200
		// metric info

		// 3. Создаём тестовый запрос (например, GET /counter/myMetric)
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/update/gauge/%s", k), nil)
		w = httptest.NewRecorder()

		// 4. Пускаем запрос через роутер
		r.ServeHTTP(w, req)

		// 5. Печатаем код ответа
		fmt.Println(w.Code)

		// 6. Печатаем тело ответа
		fmt.Println(w.Body.String())

		// Output:
		// 200
		// metric info
	}
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
