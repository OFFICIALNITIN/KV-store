package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/OFFICIALNITIN/KV-store/internal/store"
)

type HTTPHandler struct {
	kv *store.KVStore
}

func NewHTTPHandler(kv *store.KVStore) *HTTPHandler {
	return &HTTPHandler{kv: kv}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// we expect the URL to be: /keys/{key}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[1] != "keys" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	key := parts[2]

	switch r.Method {
	case http.MethodGet:

		val, found := h.kv.Get(key)
		if !found {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]any{"key": key, "value": val})

	case http.MethodPut:
		var req struct {
			Value any
			TTL   int
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		h.kv.Set(key, req.Value, 0)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))

	case http.MethodDelete:
		deleted := h.kv.Delete(key)
		if !deleted {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "deleted"}`))

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
