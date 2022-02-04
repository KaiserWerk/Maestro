package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/KaiserWerk/Maestro/internal/entity"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KaiserWerk/Maestro/internal/global"

	"github.com/gorilla/mux"
)

func TestHttpHandler_RegistrationHandler_Success(t *testing.T) {
	reg := entity.Registrant{
		Id:      "test-service",
		Address: "http://some-addr.com",
	}
	b, _ := json.Marshal(reg)
	handler := BaseHandler{}
	req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", bytes.NewBuffer(b))
	req.Header.Add("X-Registry-Token", "123abc")
	w := httptest.NewRecorder()
	handler.RegistrationHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("expected status code 200, got '%d' (%s)", resp.StatusCode, resp.Status)
	}
}

func TestHttpHandler_RegistrationHandler_Failure(t *testing.T) {
	reg := entity.Registrant{
		Id:      "test-service",
		Address: "http://some-addr.com",
	}
	b, _ := json.Marshal(reg)
	handler := BaseHandler{}
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", bytes.NewBuffer(b))
	handler.RegistrationHandler(w, req)
	req = httptest.NewRequest(http.MethodPost, "http://example.com/foo", bytes.NewBuffer(b))
	handler.RegistrationHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != 409 {
		cont, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("expected status code 409, got '%d' (%s) (%s)", resp.StatusCode, resp.Status, cont)
	}
}

func BenchmarkHttpHandler_RegistrationHandler(b *testing.B) {
	handler := &BaseHandler{
		Logger: nil, // TODO
	}

	port := global.GetPortForTest()
	router := mux.NewRouter()
	router.HandleFunc("/", handler.RegistrationHandler)

	srv := http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("localhost:%d", port),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	go func() {
		_ = srv.ListenAndServe()
	}()

	client := &http.Client{Timeout: time.Second}
	body := `{"id": "some-service-handle","address": "http://localhost:9001"}`
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/", port), bytes.NewBufferString(body))
		_, err := client.Do(req)
		if err != nil {
			b.Fatalf("could not execute request: %s", err.Error())
		}
	}

	b.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			b.Fatalf("cannot shut down HTTP server: %s", err.Error())
		}
	})
}
