package handler

import (
	"net/http"
)

func (h *BaseHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.Body.Close()
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id query parameter missing or empty", http.StatusBadRequest)
		return
	}

	err := h.MaestroCache.Update(id)
	if err != nil {
		http.Error(w, "could not update entry", http.StatusInternalServerError)
		return
	}
}
