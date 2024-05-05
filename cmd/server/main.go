/*

Разработайте сервер для сбора рантайм-метрик, который будет собирать репорты от агентов по протоколу HTTP. Агент вам предстоит реализовать в следующем инкременте — в качестве источника метрик вы будете использовать пакет runtime.
Сервер должен быть доступен по адресу http://localhost:8080, а также:
Принимать и хранить произвольные метрики двух типов:
Тип gauge, float64 — новое значение должно замещать предыдущее.
Тип counter, int64 — новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
Принимать метрики по протоколу HTTP методом POST.
Принимать данные в формате http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>, Content-Type: text/plain.
При успешном приёме возвращать http.StatusOK.
При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest.
Редиректы не поддерживаются.
Для хранения метрик объявите тип MemStorage. Рекомендуем использовать тип struct с полем-коллекцией внутри (slice или map). В будущем это позволит добавлять к объекту хранилища новые поля, например логер или мьютекс, чтобы можно было использовать их в методах. Опишите интерфейс для взаимодействия с этим хранилищем.
Пример запроса к серверу:
POST /update/counter/someMetric/527 HTTP/1.1
Host: localhost:8080
Content-Length: 0
Content-Type: text/plain
Пример ответа от сервера:
HTTP/1.1 200 OK
Date: Tue, 21 Feb 2023 02:51:35 GMT
Content-Length: 11
Content-Type: text/plain; charset=utf-8
*/

package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"metrics/internal/storage/memstorage"
)

var memStorage = memstorage.New()

func main() {
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/{typeMetric}/{nameMetric}/{valueMetric}", postMetrics)
	})
	http.ListenAndServe(":8080", r)
}

func postMetrics(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("nameMetric")
	typeMetric := r.PathValue("typeMetric")
	value := r.PathValue("valueMetric")
	if name == "" || value == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	if typeMetric != "gauge" && typeMetric != "counter" {
		w.WriteHeader(http.StatusBadRequest)
	}

	if r.PathValue("typeMetric") == "gauge" {
		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		memStorage.Set(name, valueFloat)
		w.Header().Add("Content-Type", "text-plain")
		w.Write([]byte(fmt.Sprintf(name, value)))
	}
	if r.PathValue("typeMetric") == "counter" {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		memStorage.Update(name, int64(valueInt))
		w.Header().Add("Content-Type", "text-plain")
		w.Write([]byte(fmt.Sprintf(name, value)))
	}
}
