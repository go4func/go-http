package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(err)
		}
		// log.Println("remoteAddress", r.RemoteAddr)
		resp := Resp{Header: "this is respose header", Body: "this is response body"}
		response, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(contentType, contentTypeJson)
		w.Write(response)
	}))
	defer server.Close()

	req := Req{URL: server.URL, Header: "this is request header", Body: " this is request body"}
	Send(req)
}

func TestHTTP(t *testing.T) {
	req := Req{
		URL:    "/test",
		Header: "this is  request header",
		Body:   "this is request body",
	}
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("REQUEST:", string(b))

	h := helloHandler{}
	rec := httptest.NewRecorder()
	// rec.Body = bytes.NewBuffer([]byte(""))
	request := httptest.NewRequest(http.MethodPost, req.URL, bytes.NewBuffer(b))
	// ctx := req.Context()
	// ctx = context.WithValue(ctx, "app.auth.token", "abc123")
	// ctx = context.WithValue(ctx, "app.user",
	//     &YourUser{ID: "qejqjq", Email: "user@example.com"})

	// // Add our context to the request: note that WithContext returns a copy of
	// // the request, which we must assign.
	// req = req.WithContext(ctx)

	h.ServeHTTP(rec, request)
	fmt.Println("RESPONSE:", rec.Body.String())
}

func TestSendWithContext(t *testing.T) {
	h := http.HandlerFunc(contextHandler)

	resp := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPost, "/test", nil)
	if err != nil {
		panic(err)
	}
	ctx := context.WithValue(req.Context(), "user", "nthlong")
	req = req.WithContext(ctx)

	h.ServeHTTP(resp, req)

	fmt.Println(resp.Body)
}
