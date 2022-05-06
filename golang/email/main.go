package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	http.ListenAndServe("", mux)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		resp, _ := json.Marshal(map[string]interface{}{
			"error": fmt.Sprintf("method not POST, got %s", r.Method),
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp)
		return
	}

	if r.Body == nil {
		resp, _ := json.Marshal(map[string]interface{}{
			"error": fmt.Sprintf("body is nil"),
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	defer func() {
		if _err := r.Body.Close(); _err != nil {
			log.Printf("error close request body: %s\n", _err)
			return
		}
	}()

	if err != nil {
		resp, _ := json.Marshal(map[string]interface{}{
			"error": fmt.Sprintf("failed read request body: %s", err),
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp)
		return
	}

	// only echo on success
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(reqBody)
	return
}
