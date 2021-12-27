package handler

import (
	"encoding/json"
	"net/http"
)

func (h *HttpHandler) QueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id query parameter missing or empty", http.StatusNotFound)
		return
	}
	e, ok := h.MaestroCache.Get(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		http.Error(w, "could not encode JSON", http.StatusInternalServerError)
		return
	}
}
